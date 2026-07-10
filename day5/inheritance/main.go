package main

import "fmt"

type Address struct {
	Street string
	City   string
	State  string
	Zip    string
}

func (a *Address) String() string {
	return fmt.Sprintf("Street: %s, City: %s, State: %s, Zip: %s", a.Street, a.City, a.State, a.Zip)
}

type ContactInfo struct {
	Email string
	Phone string
}

func (c *ContactInfo) String() string {
	return fmt.Sprintf("Email: %s, Phone: %s", c.Email, c.Phone)
}

// we use this very rarely
type Company struct {
	Name string
	*Address
	*ContactInfo
	BussinessType string
}

func main() {

	company := &Company{
		Name:          "Abc",
		Address:       &Address{},
		ContactInfo:   &ContactInfo{},
		BussinessType: "sample",
	}

	fmt.Println(company)

}
