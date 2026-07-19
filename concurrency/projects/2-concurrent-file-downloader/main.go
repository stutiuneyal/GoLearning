package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	destDir        = "./concurrency/projects/2-concurrent-file-downloader/downloads"
	maxConcurrency = 3
)

var downloadUrls = []string{
	"https://microsoftedge.github.io/Demos/json-dummy-data/64KB.json",
	"https://microsoftedge.github.io/Demos/json-dummy-data/128KB.json",
	"https://microsoftedge.github.io/Demos/json-dummy-data/1MB.json",
	"https://microsoftedge.github.io/Demos/json-dummy-data/10MB.json",
	"https://microsoftedge.github.io/Demos/json-dummy-data/unterminated.json",
	"https://microsoftedge.github.io/Demos/json-dummy-data/binary-data.json",
}

type Result struct {
	URL      string
	Filename string
	Size     int64
	Duration time.Duration
	Error    error
}

func DownloadFile(url string) Result {

	// extract the filename from the URL, and prepare the download destination
	fileName := filepath.Base(url)               // 1MB.json
	filePath := filepath.Join(destDir, fileName) // ./downloads/1MB.json

	result := Result{
		URL:      url,
		Filename: fileName,
	}

	fmt.Println("Downloading: ", url)

	start := time.Now()

	// Send an HTTP request
	resp, err := http.Get(url)
	if err != nil {
		result.Error = err
		result.Duration = time.Since(start)
		return result
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Println("Error closing response body: ", err)
		}
	}()

	// Only create the local file for a successfult Http response
	if resp.StatusCode != http.StatusOK {
		result.Error = fmt.Errorf("bad status: %s", resp.Status)
		result.Duration = time.Since(start)
		return result
	}

	// Create the destination file
	out, err := os.Create(filePath)
	if err != nil {
		result.Error = err
		result.Duration = time.Since(start)
		return result
	}
	defer func() {
		if err := out.Close(); err != nil {
			fmt.Println("error closing file: ", err)
		}
	}()

	// copy the response body to the file
	size, err := io.Copy(out, resp.Body)
	if err != nil {
		// remove the incomplete file
		_ = os.Remove(filePath)

		result.Error = err
		result.Duration = time.Since(start)
		return result
	}

	result.Size = size
	result.Duration = time.Since(start)

	return result

}

func ConcurrentDownloader() error {

	if maxConcurrency <= 0 {
		return fmt.Errorf("max concurrency must be greater than 0")
	}

	// create the directory if not exist
	if err := os.MkdirAll(destDir, 0o700); err != nil {
		return fmt.Errorf("failed to create destination directory: %v", err)
	}

	// Workers will send Result through results channel
	resultsChan := make(chan Result)

	// A buffered channel as a semaphore
	// At most concurrency downloads can run simultaneously
	limiter := make(chan struct{}, maxConcurrency)

	var wg sync.WaitGroup

	start := time.Now()

	for _, url := range downloadUrls {

		wg.Add(1)

		go func(downloadUrl string) {
			defer wg.Done()

			// acquire a concurrency slot
			limiter <- struct{}{}

			// Release the slot when te go routine finishes
			defer func() {
				<-limiter
			}()

			result := DownloadFile(downloadUrl)

			// Send this result back to the collector
			resultsChan <- result

		}(url)

	}

	// Wait for the goroutines to complete and close the channel
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	failureCount, successCount, totalSize := 0, 0, int64(0)
	errors := []error{}

	for result := range resultsChan {

		if result.Error != nil {
			fmt.Printf(
				"Failed: %s | Error: %v | Duration: %s\n",
				result.URL,
				result.Error,
				result.Duration,
			)

			failureCount++
			errors = append(errors, result.Error)
			continue

		}

		totalSize += result.Size
		successCount++

		fmt.Printf(
			"Downloaded %s (%d bytes) in %s\n",
			result.URL,
			result.Size,
			result.Duration,
		)
	}

	fmt.Println()
	fmt.Println("========== DOWNLOAD SUMMARY ==========")
	fmt.Println("Total files:", len(downloadUrls))
	fmt.Println("Successful:", successCount)
	fmt.Println("Failed:", failureCount)
	fmt.Println("Total bytes:", totalSize)
	fmt.Println("Maximum concurrency:", maxConcurrency)
	fmt.Println("Total duration:", time.Since(start))
	fmt.Println("======================================")

	if len(errors) > 0 {
		return fmt.Errorf("ecountered %d download errors", len(errors))
	}

	return nil

}

func main() {

	if err := ConcurrentDownloader(); err != nil {
		log.Println("Concurrent downloader completed with errors: ", err)
	}

	log.Println("Done")

}
