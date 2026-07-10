package main

import "fmt"

type Person struct {
	name string
	age  int
}

func NewPerson(name string, age int) *Person {
	return &Person{
		name: name,
		age:  age,
	}
}

func (p *Person) editName(name string) {
	p.name = name
}

func (p Person) printPerson() {
	fmt.Printf("%+v\n", p)
}

func modifyAge(age, year int) {
	age = age + year
}

func modifyAgePointer(age *int, year int) {
	*age = *age + year
}

func main() {
	age, year := 20, 5
	fmt.Println(age)
	modifyAge(age, year)
	fmt.Println(age)

	modifyAgePointer(&age, year)
	fmt.Println(age)

	person := NewPerson("Apurv", 16)
	person.printPerson()
	person.name = "Apurv Sirogi"
	person.age = 28
	person.printPerson()

	person.editName("Apurv Sirohi")
	person.printPerson()
}
