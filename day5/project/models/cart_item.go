package models

type CartItem struct {
	Product  Product
	Quantity int
}

func NewCartItem(product Product, quantity int) CartItem {
	return CartItem{
		Product:  product,
		Quantity: quantity,
	}
}

func (ci CartItem) TotalPrice() float64 {
	return ci.Product.Price * float64(ci.Quantity)
}
