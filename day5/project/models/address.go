package models

import "fmt"

type Address struct {
	Street string
	City   string
	State  string
	Zip    string
}

func NewAddress(street, city, state, zip string) *Address {
	return &Address{
		Street: street,
		City:   city,
		State:  state,
		Zip:    zip,
	}
}

func (a Address) String() string {
	return fmt.Sprintf("Street: %s, City: %s, State: %s, Zip: %s", a.Street, a.City, a.State, a.Zip)
}
