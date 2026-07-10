package models

import "fmt"

type Product struct {
	Id    int
	Name  string
	Price float64
	Stock int
}

func NewProduct(id int, name string, price float64, stock int) *Product {
	return &Product{
		Id:    id,
		Name:  name,
		Price: price,
		Stock: stock,
	}
}

func (p *Product) String() string {
	return fmt.Sprintf(
		"Product Details:\n"+
			"  Product ID : %d\n"+
			"  Name       : %s\n"+
			"  Price      : %.2f\n"+
			"  Stock      : %d",
		p.Id,
		p.Name,
		p.Price,
		p.Stock,
	)
}

func DisplayProducts(products map[int]Product) {
	fmt.Println("\n--- Available Products ---")

	for _, product := range products {
		fmt.Println(
			"ID:", product.Id,
			"| Name:", product.Name,
			"| Price:", product.Price,
			"| Stock:", product.Stock,
		)
	}
}
