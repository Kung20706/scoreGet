package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	model "test/model"
	"test/pages"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"golang.org/x/net/html"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	port          = 8707
	maxAttempts   = 5
	retryInterval = 11
	dbhost        = "127.0.0.1:3306"
)

func main() {

	// 連線到 MySQL 資料庫
	dsn := "db_user:db_password@tcp(" + dbhost + ")/db_name?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatal(err)
		}
		sqlDB.Close()
	}()

	// 自動遷移（將模型結構映射到資料庫表）
	err = db.AutoMigrate(&model.TicketNumber{}, &model.LotteryType{})
	if err != nil {
		log.Fatal(err)
	}

	opts := []selenium.ServiceOption{}

	//相對路徑的方式找出chromedriver的位置
	service, err := selenium.NewChromeDriverService("./chromedriver", port, opts...)
	if err != nil {
		panic(err)
	}
	defer service.Stop()

	caps := selenium.Capabilities{"browserName": "chrome",
		"chromeOptions": map[string]interface{}{
			"args": []string{},
		},
	}
	chromeCaps := chrome.Capabilities{
		Args: []string{
			// "--headless", // 設置無頭  正式時爬取需要使用的
			// "--proxy-server=" + proxyServerURL,
		},
	}
	caps.AddChrome(chromeCaps)
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://127.0.0.1:%d/wd/hub", port))
	if err != nil {
		panic(err)
	}
	defer wd.Quit()
	// log.Print("資訊源:https://www.lkag3.com/Issue/history?lottername=CQSSC")
	// FindScore(wd, db)
	// 間隔 5 秒
	pages.FlaCashoflotteryusa(wd, db)
	// interval := 5 * time.Second

	// // 使用 time.Tick 創建定時器
	// ticker := time.Tick(interval)
	// go func() {
	// 	for {
	// 		select {
	// 		case <-ticker:
	// 			pages.Superlotto638("", "", db)
	// 		}
	// 	}
	// }()
	// time.Sleep(2 * time.Second)
}

type LotteryResult struct {
	Period      string `json:"期別"`
	SpecialNum  string `json:"特別號"`
	WinningNums []int  `json:"獎號"`
	DrawDate    string `json:"開獎日期"`
}

// 使用模擬方式去拿資料 lkag
func FindScore(wd selenium.WebDriver, db *gorm.DB) {

	// 取得 第一個分頁的遊戲表(包括跨境遊戲)
	if err := wd.Get("https://www.lkag3.com/index/lotterylist"); err != nil {
		panic(err)
	}

	Source, err := wd.PageSource()
	if err != nil {
		log.Fatalf("Failed to get page source: %v", err)
	}

	elementTag := "href"
	elementtitle := "a"
	contains := "lottername="
	//取彩種
	lotternames := FindEleByHTML(Source, elementtitle, elementTag, contains)
	for i, lottername := range lotternames[15:35] {
		condition := model.LotteryType{
			Name: lottername,
		}
		// 找出球種 若沒有則新增到球種表
		var result model.LotteryType
		if err := db.Where(condition).FirstOrCreate(&result).Error; err != nil {
			log.Fatal(err)
		}
		for attempt := 0; attempt < maxAttempts; attempt++ {
			if err := wd.Get("https://www.lkag3.com/Issue/history?lottername=" + lottername); err != nil {
				panic(err)
			}

			pageSource, err := wd.PageSource()
			if err != nil {

				log.Fatalf("Failed to get page source: %v", err)
			}

			// 解析超文本字串
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(pageSource))
			if err != nil {
				log.Fatal(err)
			}

			var result model.LotteryType
			if err := db.Where("name = ?", lottername).First(&result).Error; err != nil {
				log.Fatal(err)
			}

			// 如果找到了，result 就包含了對應的類型
			if result.ID != 0 {
				fmt.Printf("LotteryType ID for %s: %d\n", lottername, result.ID)
			} else {
				fmt.Printf("LotteryType not found for %s\n", lottername)
			}

			td := doc.Find("td.ball")

			if td.Length() == 0 {
				fmt.Println("No matching <td> element found")
				time.Sleep(10 * time.Second)
				continue // 重新調用迴圈:用來支援請求時尚未渲染網頁 取得不了資訊的問題
			}
			// 從 <td> 內的 <span> 元素中提取內容
			time.Sleep(5 * time.Second)
			// 遍歷 1 到 10 的行數
			for i := 10; i >= 1; i-- {
				// 使用 strconv.Itoa 將數字轉換為字符串
				selector := "body > div.lotter-history-content > div > table > tbody > tr:nth-child(" + strconv.Itoa(i) + ")> td:nth-child(2)"
				selectordate := "body > div.lotter-history-content > div > table > tbody > tr:nth-child(" + strconv.Itoa(i) + ")> td:nth-child(3)"
				rowset := model.TicketNumber{}
				// 使用 Each 方法處理每個匹配到的元素
				rowset.LotteryTypeID = result.ID

				doc.Find(selector).Each(func(j int, s *goquery.Selection) {
					rowset.LotteryDay = s.Text()
				})

				doc.Find(selectordate).Each(func(j int, s *goquery.Selection) {
					rowset.StartTime = s.Text()
				})
				// Find all spans inside the div with class "b1"
				spans := doc.Find("tr:nth-child(" + strconv.Itoa(i) + ") > td:nth-child(1) > div.b4 span, tr:nth-child(" + strconv.Itoa(i) + ") > td:nth-child(1) > div.b3 span, tr:nth-child(" + strconv.Itoa(i) + ") > td:nth-child(1) > div.b2 span, tr:nth-child(" + strconv.Itoa(i) + ") > td:nth-child(1) > div.b1 span, tr:nth-child(1) > td.ball > table.custom-ball > tbody > tr.upper-row > td > span, tr:nth-child(" + strconv.Itoa(i) + ")>td.ball > div > span")

				title := doc.Find("body > div.history-navigater > a.fenlei").Each(func(j int, s *goquery.Selection) {
					log.Print(s.Text())
				})
				var result model.LotteryType
				condition := model.LotteryType{
					Name: lottername,
				}
				if err := db.Where(condition).First(&result).Error; err != nil {
					log.Fatal(err)
				}

				// 找到資料後，更新
				result.Name = lottername
				result.Namech = title.Text()
				db.Save(&result)

				var resultBuilder strings.Builder
				spans.Each(func(j int, s *goquery.Selection) {
					// 處理每個 span row
					// fmt.Printf("Row %d, Span %d: %s\n", i, j+1, s.Text())
					if j > 0 {
						resultBuilder.WriteString(",")
					}
					resultBuilder.WriteString(s.Text())
				})
				rowset.WinningNumber = resultBuilder.String()

				log.Print(rowset)

				err = db.Where(rowset).FirstOrCreate(&rowset).Error
				if err != nil {
					log.Fatal(err)
				}

			}

			// 判断是否回傳429错误
			if strings.Contains(pageSource, "429 Too Many Requests") {
				log.Println("Received 429 Too Many Requests. Waiting for a while and retrying...")
				time.Sleep(15 * time.Second)
				i-- // 重请求
				continue
			}
			// Break out of the loop if successful
			break
		}
	}
}
func FindEleByHTML(htmlString string, title string, name string, contains string) []string {
	var lotternames []string

	doc, err := html.Parse(strings.NewReader(htmlString))
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return lotternames
	}

	var findLottername func(*html.Node)
	findLottername = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == title {
			for _, attr := range n.Attr {
				if attr.Key == name && strings.Contains(attr.Val, contains) {
					// Extract lottername from the href attribute
					lottername := strings.TrimPrefix(attr.Val, "https://www.lkag3.com/Issue/history?lottername=")
					lotternames = append(lotternames, lottername)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findLottername(c)
		}
	}

	findLottername(doc)
	return lotternames
}
