package main

import (
	"fmt"
	"time"
)

func parseDate(input string) (string, error) {
	// Define the layout pattern of the input date string
	// Note: In Go, the reference time is "Mon Jan 2 15:04:05 MST 2006"
	const inputLayout = "Jan 2, 2006"

	// Define the desired output format
	const outputLayout = "2006-01-02"

	// Parse the input string into time.Time
	parsedTime, err := time.Parse(inputLayout, input)
	if err != nil {
		return "", err
	}

	// Format the time into the desired output format
	return parsedTime.Format(outputLayout), nil
}

func main() {
	inputDate := "Dec 12, 2023"
	formattedDate, err := parseDate(inputDate)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return
	}
	fmt.Println("Formatted Date:", formattedDate)
}
