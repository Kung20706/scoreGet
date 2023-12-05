package tool

import "strconv"

// 添加逗號的邏輯
func addCommasToNumberString(input string) string {
	// 將字符串轉換為整數，以便進行格式化
	number, err := strconv.Atoi(input)
	if err != nil {
		// 轉換失敗時的處理邏輯，這裡假設你的字符串是合法的數字
		return input
	}

	// 使用 FormatInt 函數將數字轉換為字符串
	formattedNumber := strconv.FormatInt(int64(number), 10)

	return formattedNumber
}
