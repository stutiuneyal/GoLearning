package main

import (
	"fmt"
	"sync"
)

type User struct {
	name string
}

func (u User) String() string {
	return fmt.Sprintf("Hello %s", u.name)
}

func main() {

	// ch := make(chan int) // un-buffered channels -> capacity = 0 -> they cannot store any value

	// // sending
	// ch <- 10
	// // receiving
	// value := <-ch

	// fmt.Println(value)

	messages := make(chan string)

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()
		fmt.Println("Sending message to the messages channel")
		messages <- "Hello World"
		messages <- "Hello New World"
	}()

	fmt.Println("About to receive messages from go routine")
	msg := <-messages
	fmt.Println(msg)

	fmt.Println("About to receive another messages from go routine")
	msg = <-messages
	fmt.Println(msg)

	fmt.Println("----------")

	wg.Add(1)

	userCh := make(chan User)
	go func() {
		defer wg.Done()
		fmt.Println("Sending message to users channel ...")
		userCh <- User{
			name: "John Doe",
		}
	}()

	fmt.Println("About to receive message from users channel ...")
	user := <-userCh
	fmt.Println(user)

	wg.Wait()

}
