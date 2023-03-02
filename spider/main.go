package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gashon/spider/scraper"
)

func generateFileName(substr string) string {
	// Get the current date in the format "2006-01-02"
	currentDate := time.Now().Format("2006-01-02")

	// Generate a random nonce as a 4-byte integer
	rand.Seed(time.Now().UnixNano())
	randNonce := rand.Intn(9999)

	// Format the nonce as a zero-padded string
	nonceStr := fmt.Sprintf("%04d", randNonce)

	// Combine the date and nonce into a filename string
	filename := fmt.Sprintf("%s.%s.%s.csv", substr, currentDate, nonceStr)

	return filename
}

func main() {
	fmt.Println("Starting scraper")

	// undergrad
	filter := []scraper.StanfordFilterPayload{
		{
			FieldUseID: "193",
			FieldType: 1,
			FieldValue: "Vice Provost for Undergraduate Education",
		},
	}
	spider := scraper.NewScrapper(generateFileName("undergrad"), filter)
	spider.Scrape()

	// grad
	filter = []scraper.StanfordFilterPayload{
		{
			FieldUseID: "192",
			FieldType: 1,
			FieldValue: "University - Student - Graduate",
		},
	}
	spider = scraper.NewScrapper(generateFileName("grad"), filter)
	spider.Scrape()

	fmt.Println("Finished scraping")
}