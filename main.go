package main


import (
    "fmt"
    "log"
    "time"
    "github.com/PuerkitoBio/goquery"
    "github.com/tebeka/selenium"
    "github.com/tebeka/selenium/chrome"
    
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
		// 正確的方式是使用 sql.DB 的 Close 方法
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatal(err)
		}
		sqlDB.Close()
	}()
    // 自動遷移（將模型結構映射到資料庫表）
    err = db.AutoMigrate(&TicketNumber{})
    if err != nil {
        log.Fatal(err)
    }
    // 查詢
    result :=TicketNumber{
        LotteryTypeID:1,
        WinningNumber:"j",

    }
    err = db.Create(&result).Error
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
    log.Print(
        result)
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

	fmt.Println("Lotternames:")
	var  ScoreList []string
	for _, lottername := range lotternames {
		fmt.Println(lottername)
        ScoreList = append(ScoreList, lottername)
        log.Print(ScoreList)
	}
    
    for i,lottername  := range ScoreList[38:]{
        for attempt := 0; attempt < maxAttempts; attempt++ {
            if err := wd.Get("https://www.lkag3.com/Issue/history?lottername="+ lottername); err != nil {
            panic(err)
        }

    log.Print(lottername)
    time.Sleep(2 * time.Second)
    pageSource, err := wd.PageSource()
    if err != nil {
        log.Fatalf("Failed to get page source: %v", err)
    }

    // 解析超文本字串
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(pageSource))
	if err != nil {
		log.Fatal(err)
	}

    td := doc.Find("td.ball")

	if td.Length() == 0 {
		fmt.Println("No matching <td> element found")
		time.Sleep(10 * time.Second)
        continue // 重新調用迴圈:用來支援請求時尚未渲染網頁 取得不了資訊的問題 
	}
    // 從 <td> 內的 <span> 元素中提取內容
	var spans []string
	td.Find("div.b1 span, div.b2 span, div.b3 span, div.b4 span, td.v1 b1, div.gbs_bg span, tbody").Each(func(i int, span *goquery.Selection) {
		spans = append(spans, span.Text())
	})

	fmt.Println("Content of <td>: ", spans)
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


type TicketNumber struct {
	ID               int     `gorm:"column:id;primary_key;type:int(20);NOT NULL;DEFAULT:0"`	
	LotteryTypeID    int      `gorm:"column:lottery_type_id;type:int(20);"`
	WinningNumber    string    `gorm:"column:winning_number;type:varchar(50);"`
	AdditionalNumber string    `gorm:"column:additional_number;"`
	LotteryDay       string    `gorm:"column:lottery day;"` // 数据库中是 date 类型，这里使用 string 类型
	StartTime        string `gorm:"column:start_time;"`   // 数据库中是 datetime 类型，这里使用 time.Time 类型
}

// TableName 指定数据库表名
func (TicketNumber) TableName() string {
	return "Ticket_Numbers"
}

// 本地 IP 地址
func getLocalIPAddress(wd selenium.WebDriver) (string, error) {
    script := `
        return new Promise(resolve => {
            const pc = new RTCPeerConnection();
            pc.createDataChannel('');
            pc.createOffer()
                .then(offer => pc.setLocalDescription(offer))
                .catch(error => console.error(error));




            pc.onicecandidate = ice => {
                if (ice && ice.candidate && ice.candidate.candidate) {
                    const ipAddress = ice.candidate.candidate.split(' ')[4];
                    resolve(ipAddress);
                    pc.onicecandidate = null;
                }
            };
        });
    `


    result, err := wd.ExecuteScript(script, nil)
    if err != nil {
        return "", err
    }


    return result.(string), nil
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
