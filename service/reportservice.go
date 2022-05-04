package service

import (
	"database/sql"
	"time"

	u "github.com/gindemo/utils"
)

type ReportDaoService interface {
	//get all categories and total sold price or error
	AllSalesByCategories() (map[string]float64, error)
	//get total sold price every day in a week or error
	SalesInWeek(time.Time) (map[string]float64, error)
	//get all products and sold count or error
	GetAllSales() (map[string]int, error)
}
type reportDaoImpl struct {
	db *sql.DB
}

func NewReportDaoService(dbCon *sql.DB) ReportDaoService {
	return &reportDaoImpl{db: dbCon}
}

//get all categories and total sold price or error
func (ms *reportDaoImpl) AllSalesByCategories() (map[string]float64, error) {
	var m = make(map[string]float64)
	//totalSalesbyCategory is a view
	r, err := ms.db.Query("select category.name,totalSalesbyCategory.sum from totalSalesbyCategory,category where categoryid=category.id")
	if err != nil {
		return map[string]float64{}, err
	}
	for r.Next() {
		cname := new(string)
		amount := new(float64)
		r.Scan(cname, amount)
		m[*cname] = *amount
	}
	return m, nil
}

//get total sold price every day in a week or error
func (mc *reportDaoImpl) SalesInWeek(date time.Time) (map[string]float64, error) {
	dates := u.GetDatesInWeek(date)
	var m = make(map[string]float64)
	r, err := mc.db.Query("select created,sum(total_price) from orders where created between $1 and $2 group by created", dates[0], dates[6])
	if err != nil {
		return map[string]float64{}, err
	}
	for r.Next() {
		d := new(time.Time)
		p := new(float64)
		r.Scan(d, p)
		m[d.Format("Mon")] = *p
	}
	return m, nil
}

//get all products and sold count or error
func (ms *reportDaoImpl) GetAllSales() (map[string]int, error) {
	var m = make(map[string]int)
	r, err := ms.db.Query("select product.name,sum(order_product.qty) from product,order_product where product.id=product_id group by product.name")
	if err != nil {
		return map[string]int{}, err
	}
	for r.Next() {
		name := new(string)
		count := new(int)
		r.Scan(name, count)
		m[*name] = *count
	}
	return m, nil
}
