package main


import (
    "fmt"
    // "os"
    "log"
    "strings"
    "time"


    "github.com/tebeka/selenium"
    "github.com/tebeka/selenium/chrome"
    "golang.org/x/net/html"
)


const (
    port = 8080
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
    chromeCaps := chrome.Capabilities{
        Args: []string{
            // "--headless", // set chrome headless
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
    pSource, err := wd.PageSource()
    if err != nil {
        log.Fatalf("Failed to get page source: %v", err)
    }
    log.Print(pSource)


    // if err := wd.Get("https://www.lkag3.com/Issue/history?lottername=CQSSC"); err != nil {
    //     panic(err)
    // }


    // pageSource, err := wd.PageSource()
    // if err != nil {
    //     log.Fatalf("Failed to get page source: %v", err)
    // }


    // // 调用提取函数
    // result, err := extractSpanText(pageSource)
    // if err != nil {
    //     fmt.Println("Error:", err)
    //     return
    // }
   
    // url := "https://www.lkag3.com/Issue/ajax_history"
	// payload := []byte("lotterycode=CQSSC&lotteryname=CQSSC")

	// // 設定 POST 請求的標頭
	// headers := map[string]interface{}{
	// 	"Accept":          "application/json, text/javascript, */*; q=0.01",
	// 	"Accept-Encoding": "gzip, deflate, br",
	// 	"Accept-Language": "zh-TW,zh;q=0.9,en-US;q=0.8,en;q=0.7",
	// 	"Cache-Control":   "no-cache",
	// 	"Content-Length":  "35",
	// 	"Content-Type":    "application/x-www-form-urlencoded; charset=UTF-8",
	// 	"Origin":          "https://www.lkag3.com",
	// 	"Pragma":          "no-cache",
	// 	"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
	// 	"X-Requested-With": "XMLHttpRequest",
	// }

	// 發送 POST 請求
    script := `const formData = new FormData();
	formData.append('lotterycode', 'CQSSC');
	formData.append('lotteryname', 'CQSSC');

var headers = new Headers();

	// 添加需要的请求头信息// 添加请求头信息
	headers.append('Accept', 'application/json, text/javascript, */*; q=0.01');
	headers.append('Accept-Encoding', 'gzip, deflate, br');
	headers.append('Accept-Language', 'zh-TW,zh;q=0.9,en-US;q=0.8,en;q=0.7');
	headers.append('Cache-Control', 'no-cache');
	headers.append('Content-Length', '35');
	headers.append('Content-Type', 'application/x-www-form-urlencoded; charset=UTF-8');
	headers.append('Origin', 'https://www.lkag3.com');
	headers.append('Pragma', 'no-cache');
	headers.append('User-Agent', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36');
	headers.append('X-Requested-With', 'XMLHttpRequest');
	// 添加其他请求头，根据需要添加更多
	
	// 创建一个包含请求头的 options 对象
	var requestOptions = {
	  method: 'POST',  // 设置请求的方法
	  headers: headers, // 将 Headers 对象传递给 headers 属性
	  body: 'lotterycode=CQSSC&lotteryname=CQSSC' // 请求的 body，可以是字符串、FormData 等
	};
	// 使用fetch发送POST请求
	fetch('https://www.lkag3.com/Issue/ajax_history', requestOptions	)
	  .then(response => response.json())
	  .then(data => {
		// 处理返回的数据
		console.log(data);
	  })
	  .catch(error => {
		console.error('Error:', error);
	  });`
      time.Sleep(5 * time.Second)
	_, err = wd.ExecuteScript(script, nil)

	if err != nil {
		log.Fatalf("Error sending POST request: %v", err)
	}
    time.Sleep(55 * time.Second)
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


func extractSpanText(htmlContent string) ([]string, error) {
    var result []string


    // 解析HTML内容
    doc, err := html.Parse(strings.NewReader(htmlContent))
    if err != nil {
        return nil, err
    }


    // 定义递归函数来提取span文本
    var extract func(*html.Node)
    extract = func(n *html.Node) {
        // 如果当前节点是span，并且包含子节点
        if n.Type == html.ElementNode && n.Data == "span" && n.FirstChild != nil {
            // 提取span内的文本
            text := strings.TrimSpace(n.FirstChild.Data)
            result = append(result, text)
        }


        // 递归处理子节点
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            extract(c)
        }
    }


    // 调用递归函数开始提取
    extract(doc)


    return result, nil
}





