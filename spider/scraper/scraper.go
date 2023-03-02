package scraper

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Scraper struct {
	auth string
}

func postRequest(conn net.Conn, req string) {
	payload := strings.NewReader(req)
	httpReq, err := http.NewRequest("POST", os.Getenv("HOST_NAME"), payload)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	httpReq.Header.Set("Content-Type", "text/plain")
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()
}

func worker(id int, requests <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	conn, err := net.Dial("tcp", os.Getenv("HOST_NAME"))
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	for req := range requests {
		postRequest(conn, req)
		fmt.Printf("Worker %d sent request: %s\n", id, req)
	}
}

func scrape() {
	var wg sync.WaitGroup

	numRequests, _ := strconv.Atoi(os.Getenv("N_REQUESTS"))
	numWorkers, _ := strconv.Atoi(os.Getenv("N_WORKERS"))

	fmt.Printf("Sending %d requests with %d workers\n", numRequests, numWorkers)

	requests := make(chan string, numRequests)

	// Populate the channel with requests
	for i := 0; i < numRequests; i++ {
		requests <- fmt.Sprintf("Request %d", i)
	}
	close(requests)

	// Launch worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i, requests, &wg)
	}

	wg.Wait()
}

func (s *Scraper) Scrape() error {
	scrape()

	return nil
}

func NewScrapper(auth string) *Scraper {
	return &Scraper{auth: auth}
}