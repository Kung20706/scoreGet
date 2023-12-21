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

func Lotto1224(backYear, backMonth string, db *gorm.DB) {
	if backYear == "" || backMonth == "" {
		backYear, backMonth = getCurrentRepublicEraAndMonth()
	}
	URL := "https://www.taiwanlottery.com.tw/lotto/Lotto1224/history.aspx"
	title := fmt.Sprintf("雙贏彩_%s_%s", backYear, backMonth)

	client := &http.Client{}
	req, err := http.NewRequest("POST", URL, nil)
	if err != nil {
		return
	}

	// Set form data
	values := req.URL.Query()
	values.Add("Lotto1224Control_history$chk", "radYM")
	values.Add("Lotto1224Control_history$dropYear", backYear)
	values.Add("Lotto1224Control_history$dropMonth", backMonth)
	values.Add("Lotto1224Control_history$btnSubmit", "查詢")
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

	COUNT_OF_1224_LOTTERY_PRIZE_NUMBER := 12
	firstNums := doc.Find(".td_w.font_black14b_center>span>span")
	var result model.TicketNumber
	dataCount := firstNums.Length() / 2 / COUNT_OF_1224_LOTTERY_PRIZE_NUMBER
	// 一局由總球數12除以24  因為有順序和大小
	for i := 0; i < dataCount; i++ {
		original_Nums := make([]string, COUNT_OF_1224_LOTTERY_PRIZE_NUMBER)

		tempSecondNums := make([]string, COUNT_OF_1224_LOTTERY_PRIZE_NUMBER)
		for j := 0; j < COUNT_OF_1224_LOTTERY_PRIZE_NUMBER; j++ {
			// 第一場  所引出 所有球的(場次 加上j數)
			tempSecondNums[j] = strings.ReplaceAll(firstNums.Eq((i*COUNT_OF_1224_LOTTERY_PRIZE_NUMBER*2)+j).Text(), " ", "")
			// 取順序排列的表  依照球數再加上12等於相同順位
			original_Nums[j] = strings.ReplaceAll(firstNums.Eq((i*COUNT_OF_1224_LOTTERY_PRIZE_NUMBER*2)+j+12).Text(), " ", "")

		}
		// var dates []string
		selector := fmt.Sprintf("#Lotto1224Control_history_dlQuery_Lotto1224_DrawTerm_%d", i)

		// log.Print(dates)
		// newTicket.LotteryDay = dates[i]

		dateselector := fmt.Sprintf("#Lotto1224Control_history_dlQuery_Lotto1224_DDate_%d", i)

		newTicket := model.TicketNumber{
			LotteryDay:      doc.Find(selector).Text(),
			StartTime:       doc.Find(dateselector).Text(),
			WinningNumber:   strings.Join(tempSecondNums, ","),
			Original_Number: strings.Join(original_Nums, ","),
			LotteryTypeID:   rspBody.ID,
		}
		db.Where(newTicket).First(&result)
		if db.Error != nil {
			fmt.Println("Failed to query records:", db.Error)
			return
		}

		// 更新记录
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
