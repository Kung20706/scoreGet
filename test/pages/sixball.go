package pages

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	model "test/model"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/tebeka/selenium"
	"gorm.io/gorm"
)

func Marksix(wd selenium.WebDriver, db *gorm.DB) {
	// 取得 第一個分頁的遊戲表(包括跨境遊戲)

	soruceurl := "https://bet.hkjc.com/marksix/Results.aspx?lang=ch###"
	if err := wd.Get(soruceurl); err != nil {
		log.Fatalf("Error opening the website: %v", err)
	}

	time.Sleep(5 * time.Second)
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
		Namech: "香港六合彩",
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
	// 內文渲染後找出html元素(球)
	roundsoruce := fmt.Sprintf("#resultMainTable > div > div")

	log.Print("doc.Find(roundsoruce)--doc.Find--", doc.Find(roundsoruce).Length())
	for j := 0; j < doc.Find(roundsoruce).Length(); j++ {
		//取doc內的彩種
		var balls string
		var sballs string
		for i := 0; i < 8; i++ {
			ballsoruce := fmt.Sprintf("#resultMainTable > div > div:nth-child(%d)>div.resultMainCell4>span:nth-child(%d)>img", j+1, i+1)

			bs := doc.Find(ballsoruce).AttrOr("src", "")
			pattern := `no_(\d+)_s\.gif`
			re := regexp.MustCompile(pattern)
			matches := re.FindStringSubmatch(bs)
			if len(matches) > 1 {

				// The first submatch (index 1) contains the captured substring
				substring := matches[1]
				fmt.Println("Extracted substring:", substring, matches)
				balls += substring
				// 號碼後面增加逗點
				if i < 6 {
					balls += ","
				}
				if i == 7 {
					sballs = matches[1]
				}
			} else {
				fmt.Println("Substring not found")
			}
			log.Print(bs)
		}

		nowandall := fmt.Sprintf("#resultMainTable > div > div:nth-child(%d) > div.resultMainCell1", j+1)
		ddmmyyyy := fmt.Sprintf("#resultMainTable > div > div:nth-child(%d) > div.resultMainCell2", j+1)

		dmy := doc.Find(ddmmyyyy).Text()

		td := doc.Find(nowandall).Text()
		newTicket := model.TicketNumber{
			WinningNumber:  balls,
			StartTime:      dmy,
			Special_Number: sballs,
			LotteryDay:     td,
			LotteryTypeID:  rspBody.ID,
			// Original_Number: historynumbersString,
			// LotteryTypeID:   rspBody.ID,
		}
		log.Print(newTicket)
		var result model.TicketNumber
		db.Where(newTicket).FirstOrCreate(&result)
		if db.Error != nil {
			fmt.Println("Failed to query records:", db.Error)
			return
		}
		log.Print(td, "ttt", dmy)
		// 如果找到符合条件的记录，更新记录
		if result.ID != 0 {

			// 使用 Update 方法更新记录
			db.Model(&model.TicketNumber{}).Where(newTicket).Updates(newTicket)
			if db.Error != nil {
				fmt.Println("Failed to update records:", db.Error)
				return
			}

			fmt.Println("Records updated successfully.")
		} else {
			fmt.Println("No records found.")
		}

	}

}
