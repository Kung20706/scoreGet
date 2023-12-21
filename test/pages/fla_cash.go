package pages

import (
	"fmt"
	"log"

	// "regexp"
	"strings"
	model "test/model"
	"test/plug"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/tebeka/selenium"
	"gorm.io/gorm"
)

func FlaCashoflotteryusa(wd selenium.WebDriver, db *gorm.DB) {
	// 取得 第一個分頁的遊戲表(包括跨境遊戲)

	soruceurl := "https://www.lotteryusa.com/florida/fantasy-5/"
	if err := wd.Get(soruceurl); err != nil {
		log.Fatalf("Error opening the website: %v", err)
	}

	time.Sleep(1 * time.Second)
	Source, err := wd.PageSource()
	if err != nil {
		log.Fatalf("Failed to get page source: %v", err)
	}
	// ballsoruce := "":nth-child(%d)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(Source))
	if err != nil {
		log.Fatal(err)
	}
	findtype := model.LotteryType{
		Namech: "佛羅里達天天樂",
	}
	var rspBody model.LotteryType
	// 使用 `FirstOrCreate` 查询符合条件的记录 新增
	if err := db.Where(findtype).FirstOrCreate(&rspBody).Error; err != nil {
		log.Fatal(err)
	}

	db.Where(findtype).First(&rspBody)
	if db.Error != nil {
		fmt.Println("Failed to query records:", db.Error)
		return
	}
	// 找出html元素(球)
	roundsoruce := fmt.Sprintf(" tr.c-result-card.c-result-card--squeeze ")
	doc.Find(roundsoruce).Each(func(i int, s *goquery.Selection) {
		// 获取日期

		date := s.Find("time.c-result-card__title").Text()
		formattedDate, err := plug.ParseDate(date)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			return
		}
		log.Print(formattedDate)
		// yyydds = append(yyydds,date)

		// 获取每个球的号码

		var resultBuilder strings.Builder
		s.Find("span.c-ball__label").Each(func(j int, span *goquery.Selection) {
			if j > 0 {
				resultBuilder.WriteString(",")
			}
			resultBuilder.WriteString(span.Text())
		})
		log.Print(resultBuilder.String())
		// newTicketNumber[i].WinningNumber=resultBuilder.String()
		// log.Print(balls,yyydds)
		fmt.Println("------")
	}).Text()

}
