package entities

import "time"

type Customer struct {
	ID      string
	Name    string
	Email   string
	Phone   string
	Address string
}

type Orders struct {
	ID               string
	Created          time.Time
	NumberOfProducts int
	Amount           float64
	Confirm          bool
}
type Customer_Order struct {
	ID        string
	Customers Customer
	Orders    Orders
}
type Order_products struct {
	Id        string
	OrderID   string
	ProductID string
	Qty       int
	Amount    float64
}

type OrderDetail struct {
	Orders
	Customer Customer
	Products []Order_ProductDetail
}
type Order_ProductDetail struct {
	Name   string
	Price  float64
	Qty    int
	Amount float64
}
