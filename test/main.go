package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"test/model"
	"strconv"
	"time"
    "github.com/PuerkitoBio/goquery"
    "github.com/tebeka/selenium"
    "github.com/tebeka/selenium/chrome"
    "strings"
    "golang.org/x/net/html"
)



const (
    port = 8707
    maxAttempts = 5
    retryInterval = 11
)

func main() {
	// 連線到 MySQL 資料庫
	dsn := "db_user:db_password@tcp(127.0.0.1:3306)/db_name?charset=utf8mb4&parseTime=True&loc=Local"
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
	err = db.AutoMigrate(&models.TicketNumber{},&models.LotteryType{})
	if err != nil {
		log.Fatal(err)
	}
	
	opts := []selenium.ServiceOption{
    }

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
    wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://127.0.0.1:%d/wd/hub",port))
    if err != nil {
        panic(err)
          }
    defer wd.Quit()

    // 取得 第一個分頁的遊戲表(包括跨境遊戲)
    if err := wd.Get("https://www.lkag3.com/index/lotterylist"); err != nil {
        panic(err)
    }

    Source, err := wd.PageSource()
    if err != nil {
        log.Fatalf("Failed to get page source: %v", err)
    }

    //找到想要元素的標籤
    elementTag:="href"
    elementtitle:="a"
    contains := "lottername="
    lotternames := FindEleByHTML(Source,elementtitle,elementTag,contains)
	var  ScoreList []string
	for _, lottername := range lotternames {
		// 查詢條件 將表的球種ID適配到爬蟲擷取的欄位
		condition:=models.LotteryType{
			Name:lottername,
		}
		// 找出球種 若沒有則新增到球種表
		var result  models.LotteryType
		if err := db.Where(condition).FirstOrCreate(&result).Error; err != nil {
			log.Fatal(err)
		}

		// 字串存儲於切片 
		ScoreList = append(ScoreList, lottername)
		}
    
    for i,lottername := range ScoreList {
        for attempt := 0; attempt < maxAttempts; attempt++ {
            if err := wd.Get("https://www.lkag3.com/Issue/history?lottername="+ lottername); err != nil {
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
	var result models.LotteryType
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
	
	// 遍歷 1 到 10 的行數
	for i := 1; i <= 10; i++ {
		// 使用 strconv.Itoa 將數字轉換為字符串
		selector := "body > div.lotter-history-content > div > table > tbody > tr:nth-child(" + strconv.Itoa(i) + ")> td:nth-child(2)"
		selectorformball := "body > div.lotter-history-content > div > table > tbody > tr:nth-child(" + strconv.Itoa(i) + ") > td.ball> div "
		selectordate := "body > div.lotter-history-content > div > table > tbody > tr:nth-child(" + strconv.Itoa(i) + ")> td:nth-child(3)"
		// 使用 Each 方法處理每個匹配到的元素
		rowset := models.TicketNumber{}
		rowset.LotteryTypeID =  result.ID 
		doc.Find(selector).Each(func(j int, s *goquery.Selection) {
			rowset.LotteryDay= s.Text()
		})
		doc.Find(selectorformball).Each(func(j int, s *goquery.Selection) {
			rowset.WinningNumber= s.Text()
		})
		doc.Find(selectordate).Each(func(j int, s *goquery.Selection) {
			rowset.StartTime= s.Text()
		})

		var spans []string
		td.Find("div.b1 span, div.b2 span, div.b3 span, div.b4 span, td.v1 b1, div.gbs_bg span").Each(func(i int, span *goquery.Selection) {
			spans = append(spans, span.Text())
			log.Print(spans,span.Text())
		})
		log.Print(rowset)
		err = db.Create(&rowset).Error
		if err != nil {
			log.Fatal(err)
		}
	
	}
  
	
    // 判断是否回傳429错误
	if strings.Contains(pageSource, "429 Too Many Requests") {
			log.Println("Received 429 Too Many Requests. Waiting for a while and retrying...")
			time.Sleep(15 * time.Second)
			i--                            // 重请求
			continue
	}
        // Break out of the loop if successful
    break
    }
}
   
    time.Sleep(55 * time.Second)
}




func FindEleByHTML(htmlString  string ,title string, name string,contains string) []string {
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
