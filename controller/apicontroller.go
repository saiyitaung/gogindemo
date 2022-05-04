package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gindemo/entities"
	"github.com/gindemo/service"
)

type ApiController struct {
	categoryService service.CategoryDaoService
	productService  service.ProductDaoService
}

func NewApiController(cs service.CategoryDaoService, ps service.ProductDaoService) *ApiController {
	return &ApiController{categoryService: cs, productService: ps}
}
func (api *ApiController) AllCategory(ctx *gin.Context) {
	categories, err := api.categoryService.FindAll()
	if err != nil {
		fmt.Println(err)
		ctx.Abort()
		return
	}
	ctx.JSON(200, categories)
}
func (api *ApiController) CreateCategory(ctx *gin.Context) {
	var category entities.Category
	err := ctx.ShouldBind(&category)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
	} else {
		api.categoryService.Create(category)
		ctx.JSON(http.StatusCreated, category)
	}

}
func (api *ApiController) GetCategoryById(ctx *gin.Context) {
	c, err := api.categoryService.FindCategoryById(ctx.Params.ByName("id"))
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx.JSON(200, c)
}
func (api *ApiController) UpdateCategory(ctx *gin.Context) {
	var c entities.Category
	err := ctx.ShouldBind(&c)
	if err != nil {
		fmt.Println(err)
		return
	}
	api.categoryService.UpdateCategory(c)
}

func (api *ApiController) AllProducts(ctx *gin.Context) {
	products, _ := api.productService.FindAll()
	ctx.JSON(200, products)
}
func (api *ApiController) GetProductById(ctx *gin.Context) {
	p, err := api.productService.FindById(ctx.Params.ByName("id"))
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx.JSON(200, p)
}
