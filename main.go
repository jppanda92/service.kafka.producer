package main

import (
	"fmt"
	"log"
	"poller"
	"scraper"
	"sync"
)

// Custom struct to hold the []string data
type StringData struct {
	Data []string `json:"data"`
}

func main() {
	companies, err := scraper.GetSP500Companies()
	if err != nil {
		log.Fatal(err)
	}

	// Writing to CSV
	err = scraper.WriteToCSV(companies, "sp500_companies.csv")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("CSV file written successfully.")

	// Reading from CSV
	readCompanies, err := poller.ReadFromCSV("sp500_companies.csv")
	if err != nil {
		log.Fatal(err)
	}

	concurrentPollers(readCompanies[60:66])

}

func concurrentPollers(symbols []string) {
	// Number of goroutines to run concurrently
	concurrentThreads := 2

	// Create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Create a channel to collect results
	results := make(chan []string, concurrentThreads)

	// Create a channel to signal errors
	errCh := make(chan error, concurrentThreads)

	// Split the symbols into chunks for each goroutine
	chunkSize := len(symbols) / concurrentThreads
	for i := 0; i < concurrentThreads; i++ {
		startIndex := i * chunkSize
		endIndex := (i + 1) * chunkSize
		if i == concurrentThreads-1 {
			endIndex = len(symbols)
		}

		// Increment the wait group for each goroutine
		wg.Add(1)

		// Launch a goroutine for each chunk of symbols
		go func(symbols []string) {
			// Decrement the wait group when the goroutine finishes
			defer wg.Done()

			// Get quotes for the symbols in this chunk
			quotesData, err := poller.GetQuotes(symbols)
			if err != nil {
				// Send the error to the error channel
				errCh <- fmt.Errorf("error getting quotes for symbols %v: %v", symbols, err)
				return
			}

			// Send the results to the channel
			results <- quotesData
		}(symbols[startIndex:endIndex])
	}

	// Close the results channel once all goroutines finish
	go func() {
		wg.Wait()
		close(results)
		close(errCh)
	}()

	// Collect results from the channel
	var allQuotesData []string
	for quotesData := range results {
		allQuotesData = append(allQuotesData, quotesData...)
	}

	// Check for errors
	if err := <-errCh; err != nil {
		log.Printf("Error occurred: %v", err)
		return
	}

	fmt.Println("All quotes data:", allQuotesData)
}
