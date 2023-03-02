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

	"github.com/gashon/spider/parser"
	"github.com/joho/godotenv"
)

type StanfordFilterPayload struct{
	FieldUseID string `json:"FieldUseID"`
	FieldType int `json:"FieldType"`
	FieldValue string `json:"FieldValue"`
}

type Scraper struct {
	Payload []StanfordFilterPayload
	FileName string
	mutex *sync.Mutex
	wg *sync.WaitGroup
}

func (s *Scraper) postRequest(conn net.Conn, req string) (string, error) {
	payload, _ := json.Marshal(s.Payload)

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



func (s *Scraper) appendResultsToFile(fileName string, profiles []parser.Person, mutex *sync.Mutex) {
	mutex.Lock()
	defer mutex.Unlock()

	file, err := os.OpenFile(fmt.Sprintf("../%s", fileName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

func (s *Scraper) worker(id int, requests <-chan string) {
	defer s.wg.Done()
	conn, err := net.Dial("tcp", os.Getenv("HOST_SOCKET"))
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	for req := range requests {
		fmt.Printf("Worker %d: Sending request %s\n", id, req)
	
		var content string
		content, err = s.postRequest(conn, req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}

		profiles := parser.Parse(content)
		s.appendResultsToFile(s.FileName, profiles, s.mutex)
	}
}

func (s *Scraper) Scrape() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return err
	}

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
		s.wg.Add(1)
		go s.worker(i, requests)
	}

	s.wg.Wait()

	return nil
}

func NewScrapper(fileName string, filter []StanfordFilterPayload) *Scraper {
	return &Scraper{
		FileName: fileName,
		Payload: filter,
		mutex: &sync.Mutex{},
		wg: &sync.WaitGroup{},
	}
}