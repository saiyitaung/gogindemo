package service

import (
	"database/sql"
	"strings"

	"github.com/gindemo/entities"
)

type ProductDaoService interface {
	//get products by category
	FindByCategoryID(string) ([]entities.Product, error)
	//get all products
	FindAll() ([]entities.Product, error)
	//create a new product
	Create(entities.Product) (entities.Product, error)
	//get product by id
	FindById(string) (entities.Product, error)
}
type productDaoServiceImpl struct {
	db *sql.DB
}

func NewProductService(dbConn *sql.DB) ProductDaoService {
	return &productDaoServiceImpl{db: dbConn}
}
func (ms *productDaoServiceImpl) FindByCategoryID(catId string) ([]entities.Product, error) {
	var products = []entities.Product{}
	r, err := ms.db.Query("select * from product where categoryId=$1", catId)
	if err != nil {
		return []entities.Product{}, err
	}
	for r.Next() {
		var p entities.Product
		var imgs string
		r.Scan(&p.ID, &p.Name, &p.Price, &p.Qty, &p.CategoryId, &imgs, &p.CoverPic, &p.Created, &p.LastUpdate, &p.Description)
		ss := strings.Split(imgs, ",")
		p.Images = ss
		products = append(products, p)
	}
	return products, nil
}
func (ms *productDaoServiceImpl) FindAll() ([]entities.Product, error) {
	var products = []entities.Product{}
	row, err := ms.db.Query("select * from product")
	if err != nil {
		return products, err
	}
	for row.Next() {
		var p entities.Product
		var imgs string
		row.Scan(&p.ID, &p.Name, &p.Price, &p.Qty, &p.CategoryId, &imgs, &p.CoverPic, &p.Created, &p.LastUpdate, &p.Description)
		ss := strings.Split(imgs, ",")
		p.Images = ss
		products = append(products, p)
	}
	return products, nil
}
func (ms *productDaoServiceImpl) Create(p entities.Product) (entities.Product, error) {
	stat := "insert into product(id,name,price,qty,categoryid,images,coverpic,created,lastupdate,description) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)"
	_, err := ms.db.Exec(stat, p.ID, p.Name, p.Price, p.Qty, p.CategoryId, strings.Join(p.Images, ","), p.CoverPic, p.Created, p.LastUpdate, p.Description)
	if err != nil {
		return p, err
	}
	return p, nil
}
func (ms *productDaoServiceImpl) FindById(id string) (entities.Product, error) {
	var p entities.Product
	r := ms.db.QueryRow("select * from product where id=$1", id)
	var imgs string
	err := r.Scan(&p.ID, &p.Name, &p.Price, &p.Qty, &p.CategoryId, &imgs, &p.CoverPic, &p.Created, &p.LastUpdate, &p.Description)
	if err != nil {
		return entities.Product{}, err
	}
	ss := strings.Split(imgs, ",")
	p.Images = ss
	return p, nil
}
