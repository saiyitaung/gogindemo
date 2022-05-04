package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/gindemo/controller"
	"github.com/gindemo/service"
)

func SetApiRoutes(r *gin.RouterGroup, dbCon *sql.DB) {
	//apis.Use(gin.BasicAuth(gin.Accounts{"someone": "somepassword"}))
	apiCtrl := controller.NewApiController(service.NewCategoryService(dbCon), service.NewProductService(dbCon))
	cApis := r.Group("/categories")
	{
		cApis.GET("", apiCtrl.AllCategory)
		cApis.GET("/:id", apiCtrl.GetCategoryById)
		// cApis.POST("", apiCtrl.CreateCategory)
		// cApis.PUT("/:id", apiCtrl.UpdateCategory)
	}
	pApis := r.Group("/products")
	{
		pApis.GET("", apiCtrl.AllProducts)
		pApis.GET("/:id", apiCtrl.GetProductById)
	}
}
