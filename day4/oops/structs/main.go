package main

import (
	"fmt"
	"strings"
	"time"
)

type Employee struct {
	Id        string
	FirstName string
	LastName  string
	Posititon string
	Salary    int
	IsActive  bool
	JoinedAt  time.Time
}

// constructor
func NewEmployee(id, firstName, lastName, position string, salary int) *Employee {
	return &Employee{
		Id:        id,
		FirstName: firstName,
		LastName:  lastName,
		Posititon: position,
		Salary:    salary,
		IsActive:  salary > 0,
		JoinedAt:  time.Now(),
	}
}

// package level functions -> encapsulation
func (e *Employee) editName(name string) error {
	names := strings.Split(strings.TrimSpace(name), " ")
	if len(names) == 1 {
		e.FirstName = names[0]
	} else if len(name) > 1 {
		e.FirstName = names[0]
		e.LastName = strings.Join(names[1:], " ")
	} else {
		return fmt.Errorf("name cannot be empty")
	}

	return nil
}

func (e *Employee) printName() {
	if e == nil {
		return
	}
	fmt.Printf("Name: %s %s\n", e.FirstName, e.LastName)
}

func main() {

	jane := NewEmployee("j123", "Jane", "Doe", "manager", 10000)

	jane.printName()

	if err := jane.editName("Jane Pamper Doe"); err != nil {
		fmt.Println(err)
		return
	}

	jane.printName()

}
