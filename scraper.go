package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/chromedp/chromedp"
)

type NLBLottery struct {
	Date       string   `json:"date"`
	DrawNumber string   `json:"draw_number"`
	Numbers    []string `json:"numbers"`
}

func main() {
	// Create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Variable to store the page content
	var htmlContent string

	// Scrape the page
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://www.nlb.lk/results/lucky-7"),
		chromedp.WaitVisible("table.tbl"), // Wait for the table to load
		chromedp.OuterHTML("html", &htmlContent),
	)
	if err != nil {
		log.Fatalf("Failed to load page: %v", err)
	}

	// Process the HTML content
	results := extractLotteryResults(htmlContent)

	// Convert results to JSON
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		log.Fatalf("Failed to convert results to JSON: %v", err)
	}

	// Print JSON to the console
	fmt.Println(string(jsonData))
}

func dsada(ticketName string) {

}

func extractLotteryResults(html string) []NLBLottery {
	var results []NLBLottery

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
				results = append(results, NLBLottery{
					Date:       date,
					DrawNumber: drawNumber,
					Numbers:    numbers,
				})
			}
		}
	}

	return results
}
