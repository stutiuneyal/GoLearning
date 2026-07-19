package main

import (
	"encoding/json"
	"fmt"
)

// var err error

// Error err = new JsonError()

// interface Animal(){}

// class Dog() implements Animal(){}

// Animal a = new Dog()

type user struct {
	Name     string `json:"name,omitempty" xml:"name"`
	Age      int    `json:"age"`
	Phone    string `json:"phone"`
	IsActive bool   `json:"is_active"`
}

func main() {

	jane := user{
		Name:     "Jane",
		Age:      32,
		Phone:    "1234567890",
		IsActive: true,
	}

	data, err := json.Marshal(jane)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("My Json Data: ", string(data))

}
