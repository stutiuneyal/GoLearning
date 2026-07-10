package main

import "fmt"

type Person interface {
	GetName() string
	PrintName()
}

type Employee struct {
	Id   string
	Name string
}

func (e *Employee) GetName() string {
	return e.Name
}

func (e *Employee) PrintName() {
	fmt.Printf("Name: %s\n", e.Name)
}

// Stringer Interface
func (e *Employee) String() string {
	return fmt.Sprintf("Id: %s, Name: %s", e.Id, e.Name)
}

type BusinessPerson struct {
	Id   string
	Name string
}

func (bp *BusinessPerson) GetName() string {
	return bp.Name
}

func (bp *BusinessPerson) PrintName() {
	fmt.Printf("Name: %s\n", bp.Name)
}

// Stringer Interface
func (bp *BusinessPerson) String() string {
	return fmt.Sprintf("Id: %s, Name: %s", bp.Id, bp.Name)
}

func displayPerson(p Person) {
	if e, ok := p.(*Employee); ok {
		fmt.Printf("Id: %s, Name: %s\n", e.Id, p.GetName())
		return
	}

	if bp, ok := p.(*BusinessPerson); ok {
		fmt.Printf("Id: %s, Name: %s\n", bp.Id, p.GetName())
		return
	}
}

func main() {

	jane := &Employee{
		Id:   "1",
		Name: "Jane Doe",
	}

	jane.PrintName()

	stuti := &BusinessPerson{
		Id:   "2",
		Name: "Stuti",
	}

	stuti.PrintName()

	displayPerson(jane)
	displayPerson(stuti)

	fmt.Println(jane)
	fmt.Println(stuti)

}
