package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

const (
	port = 8080
)

func main() {
	// 開啟 Chrome 驅動器
	opts := []selenium.ServiceOption{
		selenium.Output(os.Stderr), // 將日誌輸出到 STDERR
	}
	service, err := selenium.NewChromeDriverService("./chromedriver", port, opts...)
	if err != nil {
		log.Fatalf("Error creating ChromeDriver service: %v", err)
	}
	defer service.Stop()

	// 設定 Chrome 選項
	chromeCaps := chrome.Capabilities{
		Args: []string{
		},
	}
	caps := selenium.Capabilities{"browserName": "chrome"}
	caps.AddChrome(chromeCaps)

	// 開啟一個 Chrome 瀏覽器
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://127.0.0.1:%d/wd/hub", port))
	if err != nil {
		log.Fatalf("Error creating WebDriver: %v", err)
	}
	defer wd.Quit()

	// 設定 POST 請求的 URL 和內容
	url := "https://example.com/post-endpoint"
	payload := "key1=value1&key2=value2"

	// 使用 ExecuteScript 方法發送 POST 請求
	script := fmt.Sprintf(`
		var xhr = new XMLHttpRequest();
		xhr.open("POST", "%s", true);
		xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
		xhr.onreadystatechange = function() {
			if (xhr.readyState == 4 && xhr.status == 200) {
				console.log(xhr.responseText);
			}
		};
		xhr.send("%s");
	`, url, payload)

	_, err = wd.ExecuteScript(script, nil)
	if err != nil {
		log.Fatalf("Error executing script: %v", err)
	}

	// 等待一段時間，讓你可以看到結果
	time.Sleep(5 * time.Second)
}