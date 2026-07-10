package main

import (
	"fmt"
	"os"
)

func sampleDefer1() {
	fmt.Println("sample defer 1 executed")
}

func sampleDefer2() {
	fmt.Println("sample defer 2 executed")
}

// LIFO
func handleMultipleDefers() {
	defer sampleDefer1()
	defer sampleDefer2()

	fmt.Println("First Statement")
	fmt.Println("Second Statement")
}

func main() {

	count := 0

	defer func() {
		fmt.Println("Defer func called")
	}()

	handleMultipleDefers()

	for count < 10 {
		if count == 11 {
			return
		}
		count++
	}

	fmt.Println(count)
	fmt.Println("Main completed")
}

func sample() {
	file, err := os.Create("sample.txt")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer file.Close()

	if _, err := file.WriteString("Hello World"); err != nil {
		fmt.Println(err.Error())
		return
	}
}
