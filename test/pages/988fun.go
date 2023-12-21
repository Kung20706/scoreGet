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

func CheckScoreSec(wd selenium.WebDriver, db *gorm.DB) {
	// 取得 第一個分頁的遊戲表(包括跨境遊戲)
	if err := wd.Get("https://998.fun/html/shishicai_xj/ssc_index.html"); err != nil {
		return
	}

	time.Sleep(3 * time.Second)
	Source, err := wd.PageSource()
	if err != nil {
		log.Fatalf("Failed to get page source: %v", err)
	}
	elementTag := "tbody > tr:nth-child(2) > td.blueqiu > ul > li"

	//取doc內的彩種

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(Source))
	if err != nil {
		log.Fatal(err)
	}
	td := doc.Find(elementTag)
	log.Print(td.Text())
	STlist := model.TicketNumber{}
	td.Each(func(i int, td *goquery.Selection) {
		// 获取每个 <div> 元素的 id 属性值
		ScoreType, exists := td.Attr("id")
		if exists {
			fmt.Printf("ID 属性的值 %d: %s\n", i+1, ScoreType)
		} else {
			fmt.Printf("第 %d 个 <div> 元素没有 id 属性\n", i+1)
		}
		Score := td.Find("ul> li")
		LotteryFlag := td.Find("p")
		pattern := `第(.*?)期`
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(LotteryFlag.Text())

		if len(matches) > 1 {
			result := matches[1]
			fmt.Printf("提取的结果: %s\n", result)
			STlist.LotteryDay = result
		}
		var resultBuilder strings.Builder
		Score.Each(func(j int, li *goquery.Selection) {
			if j > 0 {
				resultBuilder.WriteString(",")
			}
			resultBuilder.WriteString(li.Text())
		})
		// 查询符合条件的记录
		condition := model.TicketNumber{
			LotteryDay:    STlist.LotteryDay,
			WinningNumber: resultBuilder.String(),
		}
		var result model.TicketNumber
		db.Where(condition).First(&result)
		if db.Error != nil {
			fmt.Println("Failed to query records:", db.Error)
			return
		}
		// 更新记录
		if result.ID != 0 {
			// 更新你需要修改的字段
			updateFields := map[string]interface{}{
				"winning_number": resultBuilder.String(),
				"lottery_day":    STlist.LotteryDay,
				"check_state":    1,
			}

			// 使用 Update 方法更新记录
			db.Model(&model.TicketNumber{}).Where(condition).Updates(updateFields)
			if db.Error != nil {
				fmt.Println("Failed to update records:", db.Error)
				return
			}

			fmt.Println("Records updated successfully.")
		} else {
			fmt.Println("No records found.")
		}

		// 输出查询结果
		fmt.Printf("Query Result: %+v\n", result)
	})

}
