package service

import (
	"database/sql"
	"fmt"

	"github.com/gindemo/entities"
)

type OrderDaoService interface {
	GetTotalCount() (int, error)
	// create a new customer's order
	CreateCustomerOrder(entities.Customer, entities.Orders, []entities.CartItem) (entities.OrderDetail, error)
	//get all orders
	FindAll(offset, limit int) ([]entities.Orders, error)
	//get order by id
	FindById(string) (entities.Orders, error)
	//Get order detail by order's id
	OrderDetail(string) (entities.OrderDetail, error)
}
type orderDaoImpl struct {
	db *sql.DB
}

func NewOrderService(dbCon *sql.DB) OrderDaoService {
	return &orderDaoImpl{db: dbCon}
}
func (ods *orderDaoImpl) createNewCustomer(c entities.Customer) (entities.Customer, error) {
	_, err := ods.db.Exec("insert into customer(id,name,email,phone,address) values($1,$2,$3,$4,$5)", c.ID, c.Name, c.Email, c.Phone, c.Address)

	if err != nil {
		return entities.Customer{}, err
	}
	return c, nil
}
func (ods *orderDaoImpl) createNewOrder(o entities.Orders) (entities.Orders, error) {
	_, err := ods.db.Exec("insert into orders(id,total_qty,total_price,confirm,created) values($1,$2,$3,$4,$5)", o.ID, o.NumberOfProducts, o.Amount, o.Confirm, o.Created)

	if err != nil {
		return entities.Orders{}, err
	}
	return o, nil
}
func (ods *orderDaoImpl) createCustomerOrder(cId, oId string) error {
	_, err := ods.db.Exec("insert into customer_order(c_id,orders_id) values($1,$2)", cId, oId)
	if err != nil {
		return err
	}
	return nil
}
func (ods *orderDaoImpl) createOrderProducts(oId, pId string, qty int, amount float64) error {
	_, err := ods.db.Exec("insert into order_product(orders_id,product_id,qty,amount) values($1,$2,$3,$4)", oId, pId, qty, amount)
	if err != nil {
		return err
	}
	return nil
}
func (ods *orderDaoImpl) GetTotalCount() (int, error) {
	r := ods.db.QueryRow("select count(*) from orders")
	var count int
	err := r.Scan(&count)
	return count, err
}
func (ods *orderDaoImpl) CreateCustomerOrder(newc entities.Customer, newo entities.Orders, cartItems []entities.CartItem) (entities.OrderDetail, error) {
	var o_detail = entities.OrderDetail{}
	c, err := ods.createNewCustomer(newc)
	if err != nil {
		fmt.Println("insert customer error : ", err)
		return o_detail, err
	}
	o, err := ods.createNewOrder(newo)
	if err != nil {
		fmt.Println("insert order error : ", err)

		return o_detail, err
	}
	err = ods.createCustomerOrder(c.ID, o.ID)
	if err != nil {
		fmt.Println("insert customer_order error : ", err)

		return o_detail, err
	}
	for _, cartItem := range cartItems {
		err = ods.createOrderProducts(o.ID, cartItem.Product.ID, cartItem.Count, cartItem.TotalAmount)
		if err != nil {
			return o_detail, err
		}
	}
	o_detail, err = ods.OrderDetail(o.ID)
	return o_detail, err
}
func (ods *orderDaoImpl) FindAll(offset, limit int) ([]entities.Orders, error) {
	var orders = []entities.Orders{}

	r, err := ods.db.Query("select * from orders offset $1 limit $2", offset, limit)
	if err != nil {
		return []entities.Orders{}, err
	}
	for r.Next() {
		var order entities.Orders
		r.Scan(&order.ID, &order.NumberOfProducts, &order.Amount, &order.Confirm, &order.Created)
		orders = append(orders, order)
	}
	return orders, nil
}
func (ods *orderDaoImpl) FindById(id string) (entities.Orders, error) {
	r := ods.db.QueryRow("select * from orders where id=$1", id)
	var o entities.Orders
	err := r.Scan(&o.ID, &o.NumberOfProducts, &o.Amount, &o.Confirm, &o.Created)
	if err != nil {
		return entities.Orders{}, err
	}
	return o, nil
}
func (ods *orderDaoImpl) OrderDetail(id string) (entities.OrderDetail, error) {
	var order_detail entities.OrderDetail
	order, err := ods.FindById(id)
	if err != nil {
		return entities.OrderDetail{}, err
	}
	order_detail.Orders = order
	r, err := ods.db.Query("select product.name,product.price,order_product.qty,order_product.amount from product,order_product where orders_id=$1 and product.id=product_id", order.ID)
	if err != nil {
		return entities.OrderDetail{}, err
	}
	cr := ods.db.QueryRow("select customer.id,name,email,phone,address from customer,customer_order where customer_order.orders_id=$1 and customer_order.c_id=customer.id", order.ID)
	var foundCustomer entities.Customer
	err2 := cr.Scan(&foundCustomer.ID, &foundCustomer.Name, &foundCustomer.Email, &foundCustomer.Phone, &foundCustomer.Address)
	if err2 != nil {
		return entities.OrderDetail{}, err
	}
	var orderProducts []entities.Order_ProductDetail
	for r.Next() {
		var op entities.Order_ProductDetail
		r.Scan(&op.Name, &op.Price, &op.Qty, &op.Amount)
		orderProducts = append(orderProducts, op)
	}
	order_detail.Products = orderProducts
	order_detail.Customer = foundCustomer
	return order_detail, nil
}
