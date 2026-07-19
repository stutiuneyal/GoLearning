package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	destDir = "./concurrency/projects/2-concurrent-file-downloader/sequential/downloads"
)

var downloadUrls = []string{
	"https://microsoftedge.github.io/Demos/json-dummy-data/64KB.json",
	"https://microsoftedge.github.io/Demos/json-dummy-data/128KB.json",
	"https://microsoftedge.github.io/Demos/json-dummy-data/1MB.json",
	"https://microsoftedge.github.io/Demos/json-dummy-data/10MB.json",
	"https://microsoftedge.github.io/Demos/json-dummy-data/unterminated.json",
	"https://microsoftedge.github.io/Demos/json-dummy-data/binary-data.json",
}

func DownloadFile(url string) error {

	// extract the filename from the URL, and prepare the download destination
	fileName := filepath.Base(url)               // 1MB.json
	filePath := filepath.Join(destDir, fileName) // ./downloads/1MB.json

	fmt.Println("Downloading: ", url)

	start := time.Now()

	// Send an HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Println("Error closing response body: ", err)
		}
	}()

	// Only create the local file for a successfult Http response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Create the destination file
	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
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

		return fmt.Errorf("failed to write file: %v", err)
	}

	fmt.Printf(
		"Downloaded %s (%d bytes) in %s\n",
		fileName,
		size,
		time.Since(start),
	)

	return nil

}

func SequentialDownloader() error {

	// create the directory if not exist
	if err := os.MkdirAll(destDir, 0o700); err != nil {
		return fmt.Errorf("failed to create destination directory: %v", err)
	}

	start := time.Now()

	successCount, failureCount := 0, 0

	for _, url := range downloadUrls {
		if err := DownloadFile(url); err != nil {
			log.Printf("Error downloading %s: %v\n\n", url, err)
			failureCount++
			continue
		}

		successCount++
		log.Println()
	}

	log.Println("========= Download Summary =========")
	log.Println("Totla Files: ", len(downloadUrls))
	log.Println("Successful: ", successCount)
	log.Println("Failed: ", failureCount)
	log.Println("Total duration: ", time.Since(start))
	log.Println("====================================")

	if failureCount > 0 {
		return errors.New(fmt.Sprintf("%d files failed to download", failureCount))
	}

	return nil
}

func main() {

	if err := SequentialDownloader(); err != nil {
		log.Println("Sequential Downloader completed with errors: ", err)
	}

	log.Println("Done")

}
