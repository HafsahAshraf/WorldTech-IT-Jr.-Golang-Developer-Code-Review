package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

type Payload struct {
	Data interface{} `json:"data"` // Change the type to interface{} to handle various types
}

func main() {
	data, err := os.ReadFile("./input.txt")
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(data), "\n")
	var wg sync.WaitGroup
	for _, line := range lines {
		if line == "" {
			continue
		}
		wg.Add(1)
		go func(line string) {
			defer wg.Done()
			data, err := getData(line)
			if err != nil {
				log.Printf("unable to get status: %v", err)
				return
			}
			if data == "foo" {
				fmt.Printf("data found: %s\n", data)
			}
		}(line)
	}
	wg.Wait()
}

func getData(line string) (string, error) {
	var location struct {
		URL string `json:"location"`
	}
	if err := json.Unmarshal([]byte(line), &location); err != nil {
		return "", err
	}

	if !strings.HasPrefix(location.URL, "http://") && !strings.HasPrefix(location.URL, "https://") {
		location.URL = "http://" + location.URL // default to http if scheme is missing
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, location.URL, nil)
	if err != nil {
		return "", err
	}

	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	// Check if the Content-Type indicates a JSON response
	contentType := res.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return "", fmt.Errorf("invalid content type: %s", contentType)
	}

	var payload struct {
		Data json.RawMessage `json:"data"` // Use json.RawMessage to handle arbitrary JSON data
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return "", err
	}

	// Return the raw JSON data
	return string(payload.Data), nil
}
