package main

import (
	"fmt"
)

func main() {

	count := 0

	if count < 0 {
		fmt.Println("Negative")
	} else if count < 20 {
		fmt.Println("Small Positive")
	} else {
		fmt.Println("Large Positive")
	}

	userAccess := map[string]bool{
		"Apurv": false,
		"Stuti": true,
	}

	if value, ok := userAccess["Apurv"]; ok { // scope of the value is only inside the if block
		fmt.Println(value)
	} else {
		fmt.Println("Key not present")
	}

	// not considered a good go practice
	access := userAccess["Apurv"]
	fmt.Println(access)

	// switch case
	day := "Saturday"

	switch day {
	case "Monday":
		fmt.Println("Monday")
	case "Tuesday", "Wednesday":
		fmt.Println("Same action on Tuesday and Wednesday")
	case "Saturday":
	case "Sunday":
		fmt.Println("Sunday")
	default:
		fmt.Println("day not in first half: ", day)
	}

	fmt.Println(count)

	/*
		Loops: for
	*/

	// simulating while in go

	// 1. infinite loop
	for {
		if count == 10 {
			break
		}
		count++
	}
	fmt.Println(count)

	count = 0

	// 2. conditional while
	for count < 10 {
		count++
	}
	fmt.Println(count)

	// 3. basic for loop
	for i := 0; i < 10; i++ {
		fmt.Print(i, " ")
	}
	fmt.Println()

	// 4. for-range -> arrays/slices, maps

	sample := []string{"a", "b", "c", "d", "e"}

	for i, v := range sample {
		fmt.Print("(", i, v, ") ")
	}
	fmt.Println()

	for _, v := range sample {
		fmt.Print("(", v, ") ")
	}
	fmt.Println()

	for i := range sample {
		fmt.Print("(", sample[i], ")")
	}
	fmt.Println()

	nameMap := map[string]int{
		"Apurv": 28,
		"Stuti": 24,
	}

	for k, v := range nameMap {
		fmt.Print("(", k, v, ") ")
	}
	fmt.Println()

	for _, v := range nameMap {
		fmt.Print("(", v, ") ")
	}
	fmt.Println()

	for k := range nameMap {
		fmt.Print("(", k, nameMap[k], ") ")
	}
	fmt.Println()

}
