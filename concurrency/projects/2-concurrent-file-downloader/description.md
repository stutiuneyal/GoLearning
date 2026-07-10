# Project Description

## Concurrent File Downloader in Go

Build a command-line file downloader that downloads files from remote URLs and stores them locally. The project should be developed in stages so that the difference between sequential and concurrent execution can be clearly observed.

### Stage 1: Download a Single File

Implement a function that accepts a file URL and downloads the file to a specified local directory.

The function must:

* Send an HTTP request using `net/http`.
* Validate the HTTP response status.
* Create a local file using `os.Create`.
* Copy the response body into the file using `io.Copy`.
* Close all resources using `defer`.
* Return an error if the request, file creation, or file-writing operation fails.

### Stage 2: Sequential File Downloader

Extend the application to download multiple files one after another.

The application must:

* Maintain a collection of file URLs.
* Create the output directory if it does not exist.
* Download each file sequentially.
* Display whether each download succeeded or failed.
* Continue downloading the remaining files even if one download fails.
* Measure and display the total execution time.

In this stage, the next download should begin only after the previous download has completed.

### Stage 3: Concurrent File Downloader

Redesign the downloader so that multiple files can be downloaded simultaneously.

The application must:

* Start each download in a separate goroutine.
* Use `sync.WaitGroup` to wait for all downloads to complete.
* Send the result of every download through a channel.
* Use a buffered channel as a semaphore to limit the number of simultaneous downloads.
* Allow the maximum concurrency limit to be configured.
* Collect and display successful and failed download results.
* Measure and display the total concurrent execution time.
* Close channels safely after all worker goroutines have completed.

Each download result should include the file URL, local file path, download status, error information, and the time taken to complete the download.

At the end of the program, display a summary containing the total number of files, successful downloads, failed downloads, and overall execution time.
