package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	start := time.Now()
	messages := make(chan string, 3)

	// fmt.Println("Sending messages to buffered channel")
	// messages <- "First Message"
	// messages <- "Second Message"
	// messages <- "Third Message"
	// messages <- "Fourth Message"

	// fmt.Println(<-messages)
	// fmt.Println(<-messages)
	// fmt.Println(<-messages)
	// fmt.Println(<-messages)

	/*
		msg := range messages <=> msg,ok := <-messages
		range -> ok handles it internally
	*/

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()
		for msg := range messages {
			time.Sleep(3 * time.Second)
			fmt.Println(msg)
		}
	}()

	/*
		for msg := range messages {} go trnaslates this to. -->

		for {
			msg,ok := <-messages
			if !ok{
				break
			}
			...
		}
	*/

	fmt.Println("Sending m1")
	messages <- "m1"

	fmt.Println("Sending m2")
	messages <- "m2"

	fmt.Println("Sending m3")
	messages <- "m3"

	fmt.Println("Sending m4 (blocked)")
	messages <- "m4"

	fmt.Println("4th sent")

	// close the channel
	close(messages)

	wg.Wait()
	fmt.Println(time.Since(start))

}
