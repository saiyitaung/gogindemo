package service

import (
	"database/sql"
	"fmt"

	"github.com/gindemo/entities"
)

type CategoryDaoService interface {
	//Get All Categories or error
	FindAll() ([]entities.Category, error)
	//Create new Cateogry
	Create(entities.Category) (entities.Category, error)
	//Get Category by Id
	FindCategoryById(id string) (entities.Category, error)
	//Update Category
	UpdateCategory(entities.Category) (entities.Category, error)
}
type CategoryServiceImpl struct {
	db *sql.DB
}

func NewCategoryService(dbCon *sql.DB) CategoryDaoService {
	return &CategoryServiceImpl{db: dbCon}
}

func (cs *CategoryServiceImpl) FindAll() ([]entities.Category, error) {
	var categories = []entities.Category{}
	row, err := cs.db.Query("select * from category")
	if err != nil {
		return []entities.Category{}, err
	}
	for row.Next() {
		var category entities.Category
		row.Scan(&category.ID, &category.Name, &category.Image)
		categories = append(categories, category)
	}
	return categories, nil
}
func (cs *CategoryServiceImpl) Create(c entities.Category) (entities.Category, error) {
	fmt.Print(c)
	_, err := cs.db.Exec("insert into category(id,name,image) values($1,$2,$3)", c.ID, c.Name, c.Image)
	if err != nil {
		return c, err
	}
	return c, nil
}
func (cs *CategoryServiceImpl) FindCategoryById(id string) (entities.Category, error) {
	r := cs.db.QueryRow("select * from category where id=$1", id)
	var c entities.Category
	err := r.Scan(&c.ID, &c.Name, &c.Image)
	if err != nil {
		return entities.Category{}, err
	}
	return c, nil
}
func (cs *CategoryServiceImpl) UpdateCategory(c entities.Category) (entities.Category, error) {
	_, err := cs.db.Exec("update category set name=$1,image=$2 where id=$3", c.Name, c.Image, c.ID)
	if err != nil {
		return entities.Category{}, err
	}
	return c, nil
}
