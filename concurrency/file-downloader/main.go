package main

import (
	"fmt"
	"math/rand"
	"time"
)

type DownloadResult struct {
	URL             string
	BytesDownloaded int
	Err             error
	Duration        int64
	Retries         int
	Seq             int
}

func main() {
	url := []string{
		"https://example.com/file1.zip",
		"https://example.com/file2.zip",
		"https://example.com/file3.zip",
		"https://example.com/file45.zip",
		"https://example.com/file3564.zip",
		"https://example.com/file3.zi456p",
		"https://example.com/file3456.zip",
		"https://example.com/file34356.zi54p",
		"https://example.com/file365546.zip",
		"https://example.com/file387697.zip",
		"https://example.com/file33245634.zip",
		"https://example.com/file3654758.zip",
		"https://example.com/file323543.zip",
		"https://example.com/file336789453.zip",
	}

	NUM_WORKER := 3
	result := make(chan DownloadResult, len(url))
	jobQueue := make(chan string)

	for i := 1; i <= NUM_WORKER; i++ {
		go func() {
			for job := range jobQueue {
				delay := time.Duration(100+rand.Intn(400)) * time.Millisecond
				time.Sleep(delay)
				result <- DownloadResult{
					URL:             job,
					BytesDownloaded: int(rand.Int31()),
					Err:             nil,
					Duration:        int64(delay.Milliseconds()),
					Retries:         1,
					Seq:             i,
				}
			}
		}()
	}

	go func() {
		for _, jobData := range url {
			jobQueue <- jobData
		}
		close(jobQueue) // signal workers no more jobs. closed for sending
	}()

	for i := 0; i < len(url); i++ {
		jobResult := <-result
		fmt.Println("Worker id :", jobResult.Seq, "Recived job data : jobResult", jobResult.URL, "downloaded in ", jobResult.Duration, "millisecond")
	}

}
