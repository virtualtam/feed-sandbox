package main

import (
	"encoding/csv"
	"log"
	"os"
)

const (
	feedsCSV = "feeds.csv"
)

func readCSV() ([]string, error) {
	file, err := os.Open(feedsCSV)
	if err != nil {
		return []string{}, err
	}
	defer file.Close()

	r := csv.NewReader(file)

	// skip the header row
	if _, err := r.Read(); err != nil {
		return []string{}, err
	}

	records, err := r.ReadAll()
	if err != nil {
		return []string{}, err
	}

	urls := make([]string, len(records))
	for i, record := range records {
		urls[i] = record[0]
	}

	return urls, nil
}

func main() {
	urls, err := readCSV()
	if err != nil {
		log.Fatal(err)
	}

	client := NewClient()

	for _, url := range urls {
		client.Fetch(url)
	}
}
