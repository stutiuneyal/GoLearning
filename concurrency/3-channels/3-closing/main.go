package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	jobs := make(chan int, 2)

	var wg sync.WaitGroup

	// producer, consumer
	wg.Add(2)

	// consumer goroutine
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		for {
			time.Sleep(3 * time.Second) // slow consumer
			r, ok := <-jobs

			if ok {
				fmt.Println("Got this message ", r)
			} else {
				fmt.Println("channel closed")
				return
			}
		}

	}(&wg)

	// producer goroutine
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		for i := 1; i <= 3; i++ {
			jobs <- i
			fmt.Println("Sending ", i)
		}

		close(jobs)
	}(&wg)

	wg.Wait()

	fmt.Println("Ending Main")
}
