package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	url = "https://microsoftedge.github.io/Demos/json-dummy-data/1MB.json"
)

func DownloadFile(destDir string) error {

	// create the directory if not exist
	if err := os.MkdirAll(destDir, 0o700); err != nil {
		return fmt.Errorf("failed to create destination directory: %v", err)
	}

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

func main() {

	if err := DownloadFile("./concurrency/projects/2-concurrent-file-downloader/downloads"); err != nil {
		log.Println("Download failed: ", err)
		return
	}

	log.Println("Done")

}
