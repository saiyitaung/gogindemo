package entities

import "fmt"

type Cart struct {
	Items []CartItem
	Total float64
}

func (c *Cart) CalcTotal() {
	var total = 0.0
	for _, item := range c.Items {
		total += item.Product.Price * float64(item.Count)
	}
	fmt.Println("total : ", total)
	c.Total = total
}

type CartItem struct {
	Product     Product
	Count       int
	TotalAmount float64
}

func (ci *CartItem) CalcPrice() {
	ci.TotalAmount = ci.Product.Price * float64(ci.Count)
}
