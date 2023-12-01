package main

import (
    "strings"
    "github.com/gocolly/colly/v2"
)

// ScrapeProxyIPs 从指定的 URL 中爬取代理 IP 地址
func ScrapeProxyIPs(url string) ([]string, error) {
    var ipAddresses []string

    // 创建一个新的收集器
    c := colly.NewCollector()

    // 设置用户代理以避免被检测为爬虫
    c.OnRequest(func(r *colly.Request) {
        r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")
    })

    // 查找和提取 IP 地址
    c.OnHTML("td", func(e *colly.HTMLElement) {
        text := strings.TrimSpace(e.Text)
        if strings.Count(text, ".") == 3 {
            ipAddresses = append(ipAddresses, text)
        }
    })

    // 访问URL
    err := c.Visit(url)
    if err != nil {
        return nil, err
    }

    return ipAddresses, nil
}