package service

import (
	"database/sql"

	"github.com/gindemo/entities"
)

type CustomerDaoService interface {
	//get toal record count
	TotalRecords() (int, error)
	//Get All Customers
	FindAll(offset, limit int) ([]entities.Customer, error)
}

type customerDaoImpl struct {
	db *sql.DB
}

func NewCustomerService(dbconn *sql.DB) CustomerDaoService {
	return &customerDaoImpl{db: dbconn}
}
func (c *customerDaoImpl) FindAll(offset, limit int) ([]entities.Customer, error) {
	var customers = []entities.Customer{}
	r, err := c.db.Query("select * from customer offset $1 limit $2", offset, limit)
	if err != nil {
		return []entities.Customer{}, err
	}
	for r.Next() {
		var customer entities.Customer
		r.Scan(&customer.ID, &customer.Name, &customer.Email, &customer.Phone, &customer.Address)
		customers = append(customers, customer)
	}
	return customers, nil
}
func (c *customerDaoImpl) TotalRecords() (int, error) {
	r := c.db.QueryRow("select count(*) from customer")
	var total int
	err := r.Scan(&total)
	return total, err
}
