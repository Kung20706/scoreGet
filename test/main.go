package main


import (
    "fmt"
    // "os"
    "log"
    "strings"
    "time"
    "golang.org/x/net/html"
    "github.com/PuerkitoBio/goquery"
    "github.com/tebeka/selenium"
    "github.com/tebeka/selenium/chrome"
)


const (
    port = 8080
	// 重做迴圈次數
    maxAttempts = 5
    retryInterval = 11
)

func main() {
    opts := []selenium.ServiceOption{
        // Enable fake XWindow session.
        // selenium.StartFrameBuffer(),


    }
    // Enable debug info.
    // selenium.SetDebug(true)
    //這裡用相對路徑的方式去寫chromedriver的位置
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
    
    // proxyServerURL := "213.157.6.50"
    chromeCaps := chrome.Capabilities{
        Args: []string{
            // "--headless", // set chrome headless
            // "--proxy-server=" + proxyServerURL,
        },
    }
    caps.AddChrome(chromeCaps)
    wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://127.0.0.1:%d/wd/hub", port))
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
	}
    log.Print(ScoreList)
    for i,lottername  := range ScoreList[49:]{
        for attempt := 0; attempt < maxAttempts; attempt++ {
            if err := wd.Get("https://www.lkag3.com/Issue/history?lottername="+ lottername); err != nil {
            panic(err)
        }

    log.Print(lottername)
    time.Sleep(3 * time.Second)
    pageSource, err := wd.PageSource()
    if err != nil {
        log.Fatalf("Failed to get page source: %v", err)
    }

    // 解析超文本字串
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(pageSource))
	if err != nil {
		log.Fatal(err)
	}

    td := doc.Find("table.dataTab")

	if td.Length() == 0 {
		fmt.Println("No matching <td> element found")
		time.Sleep(10 * time.Second)
        continue // 重新調用迴圈:用來支援請求時尚未渲染網頁 取得不了資訊的問題 
	}
    // 從 <td> 內的 <span> 元素中提取內容
	log.Print(td.Text())
	var spans []string
	
    var codelist string
	// td.Find("table.custom-ball tbody tr.upper-row").Each(func(i int, span *goquery.Selection) {
	// 	spans = append(spans, span.Text())
	// })
	td.Find("tbody tr td.issueoffoicer").Each(func(i int, span *goquery.Selection) {
		codelist, _ = td.Attr("codelist")
				
			// 解码 codelist 的 JSON 字符串
		decodedCodelist, err := unescapeCodelist(codelist)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Decoded codelist: ", decodedCodelist)

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
func unescapeCodelist(codelist string) (map[string]interface{}, error) {
    // 对 HTML 转义字符进行反转义
    unescapedCodelist, err := html.UnescapeString(codelist)
    if err != nil {
        return nil, err
    }

    // 反序列化 JSON 字符串
    var decodedCodelist map[string]interface{}
    if err := json.Unmarshal([]byte(unescapedCodelist), &decodedCodelist); err != nil {
        return nil, err
    }

    return decodedCodelist, nil
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


func unescapeCodelist(codelist string) (map[string]interface{}, error) {
    // 对 HTML 转义字符进行反转义
    unescapedCodelist, err := html.UnescapeString(codelist)
    if err != nil {
        return nil, err
    }

    // 反序列化 JSON 字符串
    var decodedCodelist map[string]interface{}
    if err := json.Unmarshal([]byte(unescapedCodelist), &decodedCodelist); err != nil {
        return nil, err
    }

    return decodedCodelist, nil
}