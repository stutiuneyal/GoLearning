package main

import (
	"fmt"
	"time"
)

func sayHello(message string, delay time.Duration) {
	time.Sleep(delay)
	fmt.Println(message)
}

func main() {

	now := time.Now()
	fmt.Println("Hello from the main() function")

	go sayHello("Hello World", time.Second)
	go sayHello("Haha", time.Second)

	// this is not the recommended way

	/*
		Why?

		1. We cannot control the order in which goroutines execute
		2. The execution time of goroutines is unpredictable, so the output may vary from one run to another
		3. If one goroutine takes 3 seconds to complete, but we only sleep for 2 seconds, the main program wil exit before the
		   that goroutine finishes, and its o/p will never be displayed
		4. We cannot keep increasing or fine-tuning the sleep duration because we don't know how long each goroutine will take,
		   in different environments or under different workloads
	*/
	time.Sleep(2 * time.Second) // Instead -> Synchronization techniques

	fmt.Println("Last hello from main() goroutine")
	fmt.Println(time.Since(now))

}
