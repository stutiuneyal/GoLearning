package models

import (
	"fmt"
	"strings"
)

type Cart struct {
	Items []CartItem
}

func NewCart() *Cart {
	return &Cart{
		Items: []CartItem{},
	}
}

/*
...v -> []v
[]v... -> ...v
*/
func (c *Cart) AddItems(items ...CartItem) {
	c.Items = append(c.Items, items...)
}

func (c *Cart) CalculateTotal() float64 {
	total := 0.0

	for _, item := range c.Items {
		total += item.TotalPrice()
	}

	return total
}

func (c *Cart) String() string {
	var result strings.Builder // StringBuilder

	result.WriteString("----- Customer Cart -----\n")

	if len(c.Items) == 0 {
		result.WriteString("Cart is empty")
		return result.String()
	}

	for _, item := range c.Items {
		result.WriteString(fmt.Sprintf(
			"%-20s | Qty: %-2d | Unit Price: ₹%.2f | Total: ₹%.2f\n",
			item.Product.Name,
			item.Quantity,
			item.Product.Price,
			item.TotalPrice(),
		))
	}

	result.WriteString("--------------------------------------------\n")
	result.WriteString(fmt.Sprintf("Cart Total: ₹%.2f", c.CalculateTotal()))

	return result.String()
}
