package main

import (
	"fmt"
)

func greet() {
	fmt.Println("Hello, World")
}

func greetName(id, name string, age int) {
	fmt.Printf("Hello, %s! Your id is: %s, and age is %d\n", name, id, age)
}

func multiply(a, b int) int {
	return a * b
}

func divide(a, b int) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("cannot divide by 0")
	}

	return float64(a) / float64(b), nil
}

/*
Variadic Function -> should be the last argument of a function -> because of this I cannot have two variadic parameters as function arguments
*/
func add(elems ...int) int {
	sum := 0
	for _, v := range elems {
		sum += v
	}
	return sum
}

// 1,2,3,4,5
// 1,2,3
// [1,2,3,4,5] -> 1,2,3,4,5
func formatFileName(fileNameFormat string, args ...interface{}) string {
	return fmt.Sprintf(fileNameFormat, args...) // Sprintf := formats and returns
}

/*
Passing slices as functions
*/
func sumOfEven(elems []int) int {
	sum := 0

	for _, v := range elems {
		if v%2 == 0 {
			sum += v
		}
	}

	return sum
}

/*
func returning a func
*/

type nextSeqFunc func(factor, offset int) int

func nextSequence(init int) nextSeqFunc {

	i := init

	return func(factor, offset int) int {
		return i*factor + offset
	}

}

func main() {

	/*
		Normal functions
	*/
	greet()
	greetName("123", "Apurv", 28)

	m := multiply(10, 20)
	fmt.Println(m)

	/*
		Error Handling
	*/
	if d, err := divide(10, 0); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(d)
	}

	/*
		Variadic Functions
	*/
	a := add(1, 2, 3)
	fmt.Println(a)
	a = add(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	fmt.Println(a)

	// fmt.Printf("fileName: %s","sample") // -> fileName: sample

	// common use-case
	fileNameFormat := "fileName: %s"
	fmt.Println(formatFileName(fileNameFormat, "sample"))

	fileNameFormat = "fileName: %s, extension: %s"
	fmt.Println(formatFileName(fileNameFormat, "sample", "pdf"))

	fileNameFormat = "fileName: %s, extension: %s, fileSize: %d kb"
	fmt.Println(formatFileName(fileNameFormat, "sample", "pdf", 200))

	/*
		passing slices as functions
	*/
	sE := sumOfEven([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	fmt.Println(sE)

	/*
		Function can return a function
	*/
	ns := nextSequence(10)
	fmt.Println(ns(10, 1))
	fmt.Println(ns(10, 11))

	/*
		anonymous functions
	*/

	func() {
		fmt.Println("Greeting from anonymous functions")
	}()

	slice1 := []int{1, 2, 3, 4}
	// slice2 := []int{2, 4, 6, 8, 10}

	mult := func(elems ...int) *int {
		mult := 1
		for _, v := range elems {
			mult *= v
		}
		return &mult
	}(slice1...)

	fmt.Println(*mult)

}
