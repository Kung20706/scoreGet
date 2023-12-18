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

func Lotto3D(backYear, backMonth string, db *gorm.DB) {
	if backYear == "" || backMonth == "" {
		backYear, backMonth = getCurrentRepublicEraAndMonth()
	}
	URL := "https://www.taiwanlottery.com.tw/Lotto/3D/history.aspx"
	title := fmt.Sprintf("三星彩_%s_%s", backYear, backMonth)

	client := &http.Client{}
	req, err := http.NewRequest("POST", URL, nil)
	if err != nil {
		return
	}

	// Set form data
	values := req.URL.Query()
	values.Add("L3DControl_history1$chk", "radYM")
	values.Add("L3DControl_history1$dropYear", backYear)
	values.Add("L3DControl_history1$dropMonth", backMonth)
	values.Add("L3DControl_history1$btnSubmit", "查詢")
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
	COUNT_OF_3D_LOTTERY_PRIZE_NUMBER := 3
	firstNums := doc.Find(".td_w.font_black14b_center")
	dataCount := firstNums.Length() / COUNT_OF_3D_LOTTERY_PRIZE_NUMBER

	for i := 0; i < dataCount; i++ {
		tempSecondNums := make([]string, COUNT_OF_3D_LOTTERY_PRIZE_NUMBER)
		for j := 0; j < COUNT_OF_3D_LOTTERY_PRIZE_NUMBER; j++ {
			tempSecondNums[j] = firstNums.Eq((i * COUNT_OF_3D_LOTTERY_PRIZE_NUMBER) + j).Text()
		}

		rowset := model.TicketNumber{}
		var dates []string
		var yyydds []string
		selector := fmt.Sprintf(" tbody > tr:nth-child(3) > td:nth-child(1)  ")
		doc.Find(selector).Each(func(i int, s *goquery.Selection) {
			date := s.Text()

			dates = append(dates, date)
		})
		// 期號:
		rowset.LotteryDay = dates[i]
		// 日期:
		doc.Find("table > tbody > tr > td > table > tbody > tr:nth-child(3) > td:nth-child(2) > p").Each(func(i int, s *goquery.Selection) {
			// 在这里处理每个符合条件的元素 s
			yyydd := s.Text()
			// 打印或处理 title
			yyydds = append(yyydds, yyydd[6:])

		})
		rowset.StartTime = yyydds[i]
		rowset.WinningNumber = strings.Join(tempSecondNums, ",")
		rowset.LotteryTypeID = rspBody.ID
		log.Print(rowset)
		// newTicket := model.TicketNumber{
		// 	WinningNumber:   numbersString,
		// 	LotteryDay:      date,
		// 	Special_Number:  doc.Find(selector).Text(),
		// 	Original_Number: tempSecondNums,
		// 	LotteryTypeID:   rspBody.ID,
		// }
		// log.Print(data, title)
	}

	time.Sleep(300 * time.Millisecond)

	return
}
