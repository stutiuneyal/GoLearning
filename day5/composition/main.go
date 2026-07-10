package main

import "fmt"

type Customer struct {
	Id              int
	Name            string
	Email           string
	BillingAddress  *Address // embedded
	ShippingAddress *Address // embedded
}

type Address struct {
	Street string
	City   string
	State  string
	Zip    string
}

func NewCustomer(id int, name, email string, billing, shipping *Address) *Customer {
	return &Customer{
		Id:              id,
		Name:            name,
		Email:           email,
		BillingAddress:  billing,
		ShippingAddress: shipping,
	}
}

// Receiver functions
func (c Customer) String() string {
	return fmt.Sprintf(`
    Customer Details:
        CustomerID: %d,
        Name: %s,
        Email: %s,
		BillingAddress: %s,
		ShippingAddress: %s
	`, c.Id, c.Name, c.Email, c.BillingAddress.String(), c.ShippingAddress.String())
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

func main() {

	customer := &Customer{
		Id:    1,
		Name:  "John",
		Email: "john@example.com",
		BillingAddress: &Address{
			Street: "abc",
			City:   "xyz",
			State:  "qwe",
			Zip:    "yuo",
		},
	}
	customer.ShippingAddress = customer.BillingAddress

	fmt.Println(customer)

}
