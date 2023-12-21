package pages

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"test/model"
	"time"

	"github.com/PuerkitoBio/goquery"
	"gorm.io/gorm"
)

func DaliyCash(backYear, backMonth string, db *gorm.DB) {
	if backYear == "" || backMonth == "" {
		backYear, backMonth = getCurrentRepublicEraAndMonth()
	}
	URL := "https://www.taiwanlottery.com.tw/lotto/DailyCash/history.aspx"
	title := fmt.Sprintf("三星彩_%s_%s", backYear, backMonth)

	client := &http.Client{}
	req, err := http.NewRequest("POST", URL, nil)
	if err != nil {
		return
	}

	// Set form data
	values := req.URL.Query()
	values.Add("D539Control_history1$chk", "radYM")
	values.Add("D539Control_history1$dropYear", backYear)
	values.Add("D539Control_history1$dropMonth", backMonth)
	values.Add("D539Control_history1$btnSubmit", "查詢")
	req.URL.RawQuery = values.Encode()

	// Perform the request
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return
	}

	lottortitle := doc.Find("#right > table > tbody > tr > td > div.font_red18.tx_md").Text()

	if doc.Find(".no_data").Length() > 0 {
		log.Printf("No data available for %s\n", title)
		return
	}
	findtype := model.LotteryType{
		Namech: strings.TrimSpace(lottortitle),
	}
	var rspBody model.LotteryType
	// 使用 `FirstOrCreate` 查询符合条件的记录 新增

	if err := db.Where(findtype).FirstOrCreate(&rspBody).Error; err != nil {
		log.Fatal(err)

	}
	COUNT_OF_539_LOTTERY_PRIZE_NUMBER := 10
	firstNums := doc.Find(".td_w.font_black14b_center")
	dataCount := firstNums.Length() / 5 / 2
	var result model.TicketNumber
	for i := 0; i < dataCount; i++ {
		tempSecondNums := make([]string, dataCount)
		for j := 0; j < COUNT_OF_539_LOTTERY_PRIZE_NUMBER; j++ {
			tempSecondNums[j] = firstNums.Eq((i * COUNT_OF_539_LOTTERY_PRIZE_NUMBER) + j).Text()
		}
		log.Print(tempSecondNums)
		// var dates []string
		selector := fmt.Sprintf("#D539Control_history1_dlQuery_D539_DrawTerm_%d", i)

		// log.Print(dates)
		// newTicket.LotteryDay = dates[i]

		dateselector := fmt.Sprintf("#D539Control_history1_dlQuery_D539_DDate_%d", i)

		newTicket := model.TicketNumber{
			StartTime:       doc.Find(dateselector).Text(),
			WinningNumber:   strings.Join(tempSecondNums, ",")[:14],
			LotteryDay:      doc.Find(selector).Text(),
			Original_Number: strings.Join(tempSecondNums, ",")[15:],
			LotteryTypeID:   rspBody.ID,
		}
		db.Where(newTicket).First(&result)
		if db.Error != nil {
			fmt.Println("Failed to query records:", db.Error)
			return
		}
		if result.ID != 0 {
			// 更新你需要修改的字段
			log.Print(newTicket)
			newTicket.CheckState = 1
			// 使用 Update 方法更新记录
			db.Model(&model.TicketNumber{}).Where(newTicket).Updates(newTicket)
			if db.Error != nil {
				fmt.Println("Failed to update records:", db.Error)
				return
			}

			fmt.Println("Records updated successfully.")
		} else {
			db.Save(&newTicket)
			fmt.Println("No records found.")
		}

	}

	time.Sleep(300 * time.Millisecond)

	return
}
