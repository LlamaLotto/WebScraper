package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
	web_scrape2 "web-scape-go/web-scrape"

	"github.com/chromedp/chromedp"
)

// Ticket names and base URL
var ticketNames = []string{
	"govisetha",
	"mahajana-sampatha",
	"dhana-nidhanaya",
	"mega-power",
	"lucky-7",
	"handahana",
}

const baseURL = "https://www.nlb.lk/results/"

// Combined structure to store all ticket results
type CombinedResults struct {
	Tickets map[string]interface{} `json:"tickets"`
}

// scrapeTicket scrapes the results for a single ticket
func scrapeTicket(ticketName string, wg *sync.WaitGroup, resultsChan chan<- map[string]interface{}) {
	defer wg.Done() // Decrement the wait group counter when done

	// Set a timeout for the scrape
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Create a new ChromeDP context
	chromeCtx, chromeCancel := chromedp.NewContext(ctx)
	defer chromeCancel()

	// Construct the URL
	url := baseURL + ticketName

	// Variable to store the page content
	var htmlContent string

	// Scrape the page
	err := chromedp.Run(chromeCtx,
		chromedp.Navigate(url),
		chromedp.WaitVisible("table.tbl", chromedp.ByQuery), // Wait for the table to load
		chromedp.OuterHTML("html", &htmlContent),
	)
	if err != nil {
		log.Printf("Failed to scrape %s: %v", ticketName, err)
		return
	}

	// Process the HTML content
	results := web_scrape2.NlbLotteryResults(htmlContent)

	// Send the results to the channel
	resultsChan <- map[string]interface{}{ticketName: results}

	log.Printf("Successfully scraped results for %s", ticketName)
}

func main() {
	startTime := time.Now() // Capture start time

	// Create output directory
	err := os.MkdirAll("results", 0755)
	if err != nil {
		log.Fatalf("Failed to create results directory: %v", err)
	}

	// Combined results structure
	combinedResults := CombinedResults{
		Tickets: make(map[string]interface{}),
	}

	// Use a wait group and a channel to manage concurrency and collect results
	var wg sync.WaitGroup
	resultsChan := make(chan map[string]interface{}, len(ticketNames))

	// Iterate over ticket names and scrape concurrently
	for _, ticketName := range ticketNames {
		wg.Add(1)
		go scrapeTicket(ticketName, &wg, resultsChan)
	}

	// Close the channel after all goroutines are done
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Collect results from the channel
	for result := range resultsChan {
		for ticketName, ticketData := range result {
			combinedResults.Tickets[ticketName] = ticketData
		}
	}

	// Save combined results to a single JSON file
	combinedFilePath := filepath.Join("results", "all_tickets.json")
	combinedJSON, err := json.MarshalIndent(combinedResults, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal combined JSON: %v", err)
	}

	err = os.WriteFile(combinedFilePath, combinedJSON, 0644)
	if err != nil {
		log.Fatalf("Failed to write combined JSON to file: %v", err)
	}

	log.Printf("Combined results saved to %s", combinedFilePath)

	endTime := time.Now()             // Capture end time
	elapsed := endTime.Sub(startTime) // Calculate duration
	fmt.Printf("Execution time: %s\n", elapsed)
}
