package plug

import (
	"fmt"
	"strings"
	"time"
)

var monthDict = map[string]string{
	"Jan": "01",
	"Feb": "02",
	"Mar": "03",
	"Apr": "04",
	"May": "05",
	"Jun": "06",
	"Jul": "07",
	"Aug": "08",
	"Sep": "09",
	"Oct": "10",
	"Nov": "11",
	"Dec": "12",
}

func ExtractMonthFromDateString(dateStr string) (string, error) {
	parts := strings.Split(dateStr, " ")
	if len(parts) < 3 {
		return "", fmt.Errorf("invalid date format")
	}

	monthAbbr := parts[1] // 提取月份縮寫，例如 "Dec"
	monthCode, ok := monthDict[monthAbbr]
	if !ok {
		return "", fmt.Errorf("invalid month abbreviation")
	}

	return monthCode, nil
}

func ParseDate(input string) (string, error) {
	// Define the layout pattern of the input date string
	// Note: In Go, the reference time is "Mon Jan 2 15:04:05 MST 2006"
	const inputLayout = "Monday, Jan 2, 2006"

	// Define the desired output format
	const outputLayout = "20060102"

	// Parse the input string into time.Time
	parsedTime, err := time.Parse(inputLayout, input)
	if err != nil {
		return "", err
	}

	// Format the time into the desired output format
	return parsedTime.Format(outputLayout), nil
}
