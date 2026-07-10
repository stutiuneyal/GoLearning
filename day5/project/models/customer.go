package models

import (
	"fmt"
	"strings"
)

type Customer struct {
	Id              int
	Name            string
	Email           string
	BillingAddress  Address
	ShippingAddress Address
}

func NewCustomer(id int, name, email string, billingAddress, shippingAddress Address) *Customer {
	return &Customer{
		Id:              id,
		Name:            name,
		Email:           email,
		BillingAddress:  billingAddress,
		ShippingAddress: shippingAddress,
	}
}

// Receiver function
func (c *Customer) UpdateEmail(email string) error {
	if !strings.Contains(email, "@") {
		return fmt.Errorf("invalid email")
	}

	c.Email = email

	return nil
}

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
