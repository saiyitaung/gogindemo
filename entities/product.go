package entities

import "time"

type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	Qty         int       `json:"qty"`
	CategoryId  string    `json:"categoryId"`
	Images      []string  `json:"images"`
	Description string    `json:"desc"`
	CoverPic    string    `json:"coverpic"`
	Created     time.Time `json:"-"`
	LastUpdate  time.Time `json:"lastaccess"`
}

func (p Product) String() string {
	return p.Name
}
