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

func Superlotto638(backYear, backMonth string, db *gorm.DB) {
	if backYear == "" || backMonth == "" {
		backYear, backMonth = getCurrentRepublicEraAndMonth()
	}
	URL := "https://www.taiwanlottery.com.tw/lotto/superlotto638/history.aspx"
	title := fmt.Sprintf("威力彩_%s_%s", backYear, backMonth)

	client := &http.Client{}
	req, err := http.NewRequest("POST", URL, nil)
	if err != nil {
		return
	}

	// Set form data
	values := req.URL.Query()
	values.Add("SuperLotto638Control_history1$chk", "radYM")
	values.Add("SuperLotto638Control_history1$dropYear", backYear)
	values.Add("SuperLotto638Control_history1$dropMonth", backMonth)
	values.Add("SuperLotto638Control_history1$btnSubmit", "查詢")
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
	// 先取各種黑紅球
	secondNums := doc.Find(".td_w.font_red14b_center")
	dataCount := secondNums.Length() / 2

	time.Sleep(300 * time.Millisecond)
	log.Print("dataCount:", dataCount)
	// "Lotto649Control_history_dlQuery_SNo_0"
	// _0->_9 所有場次的特別號
	// "Lotto649Control_history_dlQuery_No1_0"
	// 0_0為第一局 第一顆球 1_0第一局第一顆球 1_1代表第二局第一顆球
	// make a slice save struct single ball 1-6 superball 7 and  Arrange in order

	newTicketNumber := []model.TicketNumber{}
	var numbersString string
	var historynumbersString string
	// 處理每个球
	idxround := 10
	for round := 0; round < idxround; round++ {
		// 處理單一數字
		// 重置字符串
		numbersString = ""
		historynumbersString = ""
		for order := 1; order <= 6; order++ {
			selector := fmt.Sprintf("#SuperLotto638Control_history1_dlQuery_No%d_%d", order, round)
			number := doc.Find(selector).Text()

			// 添加号码到字符串
			numbersString += number

			// 在非最后一个号码后添加逗号
			if order < 6 {
				numbersString += ","
			}
			// 開球順序
			historyselector := fmt.Sprintf("#SuperLotto638Control_history1_dlQuery_SNo%d_%d", order, round)

			historynumber := doc.Find(historyselector).Text()
			// 為字符串增加號碼
			historynumbersString += historynumber
			// 號碼後面增加逗點
			if order < 6 {
				historynumbersString += ","
			}
		}
		// 取場次特别号码
		selector := fmt.Sprintf("#SuperLotto638Control_history1_dlQuery_SNo7_%d", round)
		// 取每個場次的期號
		drawselector := fmt.Sprintf("#SuperLotto638Control_history1_dlQuery_DrawTerm_%d", round)

		newTicket := model.TicketNumber{
			WinningNumber:   numbersString,
			LotteryDay:      doc.Find(drawselector).Text(),
			Special_Number:  doc.Find(selector).Text(),
			Original_Number: historynumbersString,
			LotteryTypeID:   rspBody.ID,
		}
		newTicketNumber = append(newTicketNumber, newTicket)
		log.Print("取場次特别号码", doc.Find(selector).Text(), "旗號", doc.Find(drawselector).Text())
		var result model.TicketNumber
		db.Where(newTicket).First(&result)
		if db.Error != nil {
			fmt.Println("Failed to query records:", db.Error)
			return
		}
		// 更新记录
		if result.ID != 0 {
			// 更新你需要修改的字段

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

	log.Print(newTicketNumber)

	return
}
