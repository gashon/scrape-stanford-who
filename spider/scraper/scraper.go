package scraper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
	"math/rand"

	"github.com/gashon/spider/parser"
	"github.com/joho/godotenv"
)

type Scraper struct {}

func postRequest(conn net.Conn, req string) (string, error) {
	data := []string{}
	payload, _ := json.Marshal(data)

	URL := fmt.Sprintf("%s?format=json&locale=en-US&pageNumber=%s", os.Getenv("HOST_NAME"), req)
	// URL := fmt.Sprintf("%s?format=json&locale=en-US&pageNumber=2", os.Getenv("HOST_NAME"))
	httpReq, err := http.NewRequest("POST", URL, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}
	
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return "", err
	}

	// read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return "", err
	}

	defer resp.Body.Close()

	return string(body), nil
}

func generateFileName() string {
	    // Get the current date in the format "2006-01-02"
		currentDate := time.Now().Format("2006-01-02")

		// Generate a random nonce as a 4-byte integer
		rand.Seed(time.Now().UnixNano())
		randNonce := rand.Intn(9999)
	
		// Format the nonce as a zero-padded string
		nonceStr := fmt.Sprintf("%04d", randNonce)
	
		// Combine the date and nonce into a filename string
		filename := fmt.Sprintf("%s.%s.csv", currentDate, nonceStr)
	
		return filename
}

func appendResultsToFile(profiles []parser.Person, mutex *sync.Mutex) {
	mutex.Lock()
	defer mutex.Unlock()

	file, err := os.OpenFile(fmt.Sprintf("../%s", generateFileName()), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	// put profiles in csv format
	for _, profile := range profiles {
		_, err := file.WriteString(fmt.Sprintf("%s,%s,%s,%s,%s,%s\n", profile.Name, profile.Email, profile.Affiliation, profile.Department, profile.Role, profile.ProfileURL))
		if err != nil {
			log.Fatal(err)
		}
	}
		
}

func worker(id int, requests <-chan string, mutex *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	conn, err := net.Dial("tcp", os.Getenv("HOST_SOCKET"))
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	for req := range requests {
		fmt.Printf("Worker %d: Sending request %s\n", id, req)
	
		var content string
		content, err = postRequest(conn, req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}

		profiles := parser.Parse(content)
		appendResultsToFile(profiles, mutex)
	}
}

func (s *Scraper) Scrape() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return err
	}

	var wg sync.WaitGroup
	var mutex sync.Mutex

	numRequests, _ := strconv.Atoi(os.Getenv("N_REQUESTS"))
	numWorkers, _ := strconv.Atoi(os.Getenv("N_WORKERS"))

	fmt.Printf("Sending %d requests with %d workers\n", numRequests, numWorkers)

	requests := make(chan string, numRequests)

	// Populate the channel with requests
	for i := 1; i <= numRequests; i++ {
		requests <- fmt.Sprintf("%d", i)
	}
	close(requests)

	// Launch worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i, requests, &mutex, &wg)
	}

	wg.Wait()

	return nil
}

func NewScrapper() *Scraper {
	return &Scraper{}
}