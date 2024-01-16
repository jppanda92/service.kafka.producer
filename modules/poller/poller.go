package poller

import (
	"encoding/csv"
	"os"

	"github.com/markcheno/go-quote"
)

// ReadFromCSV reads data from a CSV file and returns a slice of strings.
func ReadFromCSV(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.LazyQuotes = true
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var data []string
	for _, record := range records {
		if len(record) > 0 {
			data = append(data, record[0])
		}
	}

	return data, nil
}

// GetQuotes retrieves quotes for a list of symbols and returns the CSV data.
func GetQuotes(symbols []string) ([]string, error) {
	var quotesData []string

	for _, symbol := range symbols {
		q, err := quote.NewQuoteFromYahoo(symbol, "2016-01-01", "2016-04-01", quote.Daily, true)
		if err != nil {
			return nil, err
		}

		quotesData = append(quotesData, string(q.JSON(false)))
	}

	return quotesData, nil
}
