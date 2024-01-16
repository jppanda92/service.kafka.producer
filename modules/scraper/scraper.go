package scraper

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const sp500URL = "https://en.wikipedia.org/wiki/List_of_S%26P_500_companies"

func GetSP500Companies() ([]string, error) {
	resp, err := http.Get(sp500URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch the page: %s", resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var companies []string
	doc.Find("#constituents tbody tr td:nth-child(1)").Each(func(i int, s *goquery.Selection) {
		companyName := strings.TrimSpace(s.Text())
		companies = append(companies, companyName)
	})

	return companies, nil
}

func WriteToCSV(companies []string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, company := range companies {
		err := writer.Write([]string{company})
		if err != nil {
			return err
		}
	}

	return nil
}
