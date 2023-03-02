package main

import (
	"fmt"

	"github.com/gashon/spider/scraper"
)

func main() {

	fmt.Println("Starting scraper")

	scraper := scraper.Scraper{}

	scraper.Scrape()

	fmt.Println("Finished scraping")
}