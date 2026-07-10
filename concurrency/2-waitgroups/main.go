package main

import (
	"fmt"
	"sync"
)

func printMessage(msg string, wg *sync.WaitGroup) {
	defer wg.Done()
	if msg == "" {
		return
	}
	fmt.Println(msg)
}

func main() {

	var wg sync.WaitGroup

	messages := []string{"Hello", "World", ""}

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go printMessage(messages[i], &wg)
	}
	// wg.Add(3)

	// go printMessage("Hello", &wg)
	// go printMessage("World", &wg)
	// go printMessage("", &wg)

	wg.Wait()

	fmt.Println("All go routines completed")
}
