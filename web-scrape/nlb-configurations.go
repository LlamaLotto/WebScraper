package web_scrape

import (
	"strings"
	"web-scape-go/Model"
)

func NlbLotteryResults(html string) []Model.NLBLottery {
	var results []Model.NLBLottery

	// Parse the HTML content to extract relevant data
	rows := strings.Split(html, "<tr>")
	for _, row := range rows {
		if strings.Contains(row, "<b>") && strings.Contains(row, "<ol") {
			// Extract Draw Number
			drawStart := strings.Index(row, "<b>")
			drawEnd := strings.Index(row, "</b>")
			if drawStart == -1 || drawEnd == -1 || drawEnd < drawStart {
				continue // Skip if indices are invalid
			}
			drawNumber := strings.TrimSpace(row[drawStart+3 : drawEnd])

			// Extract Date
			dateStart := strings.Index(row, "<br>")
			if dateStart == -1 {
				continue // Skip if <br> tag is not found
			}
			dateEnd := strings.Index(row[dateStart:], "</td>")
			if dateEnd == -1 {
				continue // Skip if </td> tag is not found
			}
			date := strings.TrimSpace(row[dateStart+4 : dateStart+dateEnd])

			// Extract Numbers
			numbersStart := strings.Index(row, "<ol")
			numbersEnd := strings.Index(row[numbersStart:], "</ol>")
			if numbersStart == -1 || numbersEnd == -1 {
				continue // Skip if <ol> or </ol> tags are missing
			}
			numbersBlock := row[numbersStart : numbersStart+numbersEnd]

			// Extract individual numbers
			var numbers []string
			for _, numberBlock := range strings.Split(numbersBlock, "<li") {
				if strings.Contains(numberBlock, ">") {
					numberStart := strings.Index(numberBlock, ">")
					numberEnd := strings.Index(numberBlock[numberStart:], "<")
					if numberStart == -1 || numberEnd == -1 {
						continue // Skip if indices are invalid
					}
					number := strings.TrimSpace(numberBlock[numberStart+1 : numberStart+numberEnd])
					numbers = append(numbers, number)
				}
			}

			// Append to results
			if len(numbers) > 0 {
				results = append(results, Model.NLBLottery{
					Date:       date,
					DrawNumber: drawNumber,
					Numbers:    numbers,
				})
			}
		}
	}

	return results
}
