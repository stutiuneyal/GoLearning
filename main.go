package main

import "fmt"

type user interface {
	getuser() string
	createUser(name string) string
}

type Person struct {
	Name string
}

func (p *Person) getuser() string {
	return p.Name
}

func (p *Person) createUser(name string) string {
	p.Name = name
	return p.Name
}

// var and const block -> global

var (
	s1 string
	s2 string
)

const (
	a1 int = 32
	a2     = 64
)

func main() {

	/*
		values in go
	*/
	a := "Hello"
	b := 32    // int , int32, int64
	c := 23.45 // float32, float64
	d := true
	fmt.Println(a, b, c, d)

	// all acommodation
	var e interface{}
	e = 23.45
	fmt.Println(e)
	e = "Hello"
	fmt.Println(e)

	var f any
	f = 23.45
	fmt.Println(f)
	f = "Hello"
	fmt.Println(f)

	// Arrays
	var names [2]string
	names[0] = "Apurv"
	names[1] = "Stuti"
	fmt.Printf("%#v\n", names)

	newNames := [2]string{
		"Apurv",
		"Stuti",
	}
	fmt.Printf("%#v\n", newNames)

	// Slices
	var fullNames []string
	fullNames = append(fullNames, "Apurv")
	fullNames = append(fullNames, "Stuti")
	fmt.Printf("%#v\n", fullNames)

	newFullNames := []string{
		"Apurv",
		"Stuti",
	}
	fmt.Printf("%#v\n", newFullNames)

	makeFullNames := make([]int, 0, 10)
	fmt.Printf("%#v\n", makeFullNames)

	// Maps
	// var nameMap map[string]string{
	// 	"Apurv":"Sample",
	// }
	// nameMap["Apurv"] = "Hello"

	nameMap := map[string]string{
		"Apurv": "Hello",
		"Stuti": "World",
	}
	nameMap["John"] = "Doe"
	fmt.Printf("%#v\n", nameMap)

	fullNameMap := make(map[string]string)
	fullNameMap["John"] = "Doe"
	fmt.Printf("%#v\n", fullNameMap)

	/*
		How we can declare variables in go
	*/
	var greeting string
	greeting = "hello"
	fmt.Println(greeting)

	var year = 2025
	fmt.Println(year)

	fName := "Apurv" // var fName = "Apurv"
	fmt.Println(fName)

	newFirstName, newLastName, age := "Apurv", "Sirohi", 28
	fmt.Println(newFirstName, newLastName, age)

	/*
	 constants
	*/

	// const sample string
	// sample = "Hello"

	// sample := "hello"

	const sample = "Hello"
	// sample = "World" --> Not permissible

}
