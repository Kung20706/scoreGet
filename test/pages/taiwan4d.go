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

func Lotto4D(backYear, backMonth string, db *gorm.DB) {
	if backYear == "" || backMonth == "" {
		backYear, backMonth = getCurrentRepublicEraAndMonth()
	}
	URL := "https://www.taiwanlottery.com.tw/Lotto/4D/history.aspx"
	title := fmt.Sprintf("四星彩_%s_%s", backYear, backMonth)

	client := &http.Client{}
	req, err := http.NewRequest("POST", URL, nil)
	if err != nil {
		return
	}

	// Set form data
	values := req.URL.Query()
	values.Add("L4DControl_history1$chk", "radYM")
	values.Add("L4DControl_history1$dropYear", backYear)
	values.Add("L4DControl_history1$dropMonth", backMonth)
	values.Add("L4DControl_history1$btnSubmit", "查詢")
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

	COUNT_OF_4D_LOTTERY_PRIZE_NUMBER := 4

	firstNums := doc.Find(".td_w.font_black14b_center>span.td_w")
	length := firstNums.Length() / COUNT_OF_4D_LOTTERY_PRIZE_NUMBER
	var dates []string
	selector := fmt.Sprintf(" tbody > tr:nth-child(3) > td:nth-child(1)  ")

	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		date := s.Text()
		dates = append(dates, date)
	})
	// Ensure that index i is within the bounds of the dates slice
	rowsetlist := []model.TicketNumber{}
	rowset := model.TicketNumber{}
	var result model.TicketNumber
	for i := 0; i < length; i++ {
		tempSecondNums := make([]string, COUNT_OF_4D_LOTTERY_PRIZE_NUMBER)
		for j := 0; j < COUNT_OF_4D_LOTTERY_PRIZE_NUMBER; j++ {
			tempSecondNums[j] = firstNums.Eq((i * COUNT_OF_4D_LOTTERY_PRIZE_NUMBER) + j).Text()
		}
		rowset.WinningNumber = strings.Join(tempSecondNums, ",")
		log.Print(strings.Join(tempSecondNums, ","))
		rowset.LotteryTypeID = rspBody.ID
		if i < len(dates) {
			rowset.LotteryDay = dates[i]
		}
		rowsetlist = append(rowsetlist, rowset)
		db.Where(rowset).First(&result)
		// if db.Error != nil {
		// 	fmt.Println("Failed to query records:", db.Error)
		// 	return
		// }
		// 更新记录

		if result.ID != 0 {
			// 更新你需要修改的字段

			// 使用 Update 方法更新记录
			db.Model(&model.TicketNumber{}).Where(rowset).Updates(rowset)
			if db.Error != nil {
				fmt.Println("Failed to update records:", db.Error)
				return
			}

			fmt.Println("Records updated successfully.")
		} else {
			db.Save(&rowsetlist)
			fmt.Println("No records found.")
		}
	}

	time.Sleep(300 * time.Millisecond)
	return
}
