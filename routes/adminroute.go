package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/gindemo/controller"
	"github.com/gindemo/middlewares"
	"github.com/gindemo/service"
)

func SetAdminRoutes(g *gin.RouterGroup, dbCon *sql.DB) {

	categoryService := service.NewCategoryService(dbCon)
	orderS := service.NewOrderService(dbCon)
	customerS := service.NewCustomerService(dbCon)
	rps := service.NewReportDaoService(dbCon)
	ps := service.NewProductService(dbCon)

	adminCtrl := controller.NewAdminCtrl(categoryService, orderS, customerS, rps, ps)
	g.Use(middlewares.IsAuthenticated())
	//g.Use(gin.BasicAuth(gin.Accounts{"user": "admin"}))
	{
		{
			g.GET("", adminCtrl.AdminHome)
			g.POST("/upload", adminCtrl.UploadPic)
			g.GET("/logout", adminCtrl.Logout)
			g.GET("/error", adminCtrl.Error)
		}
		cgroutes := g.Group(controller.ADMIN_CATEGORY_ROUTE)
		{
			cgroutes.GET("", adminCtrl.AdminCategory)
			cgroutes.POST("", adminCtrl.CreateCategory)
			cgroutes.GET("/new", adminCtrl.AdminCategoryNewForm)
		}
		prodroutes := g.Group(controller.ADMIN_PRODUCTS_ROUTE)
		{
			prodroutes.GET("", adminCtrl.AdminProducts)
			prodroutes.POST("", adminCtrl.AdminProductsPost)
			prodroutes.GET("/new", adminCtrl.AdminProductsNewForm)
		}
		customers := g.Group(controller.ADMIN_CUSTOMERS_ROUTE)
		{
			customers.GET("", adminCtrl.Customers)
		}
		orders := g.Group(controller.ADMIN_ORDERS_ROUTE)
		{
			orders.GET("", adminCtrl.GetOrders)
			orders.GET("/:id", adminCtrl.GetOrderDetail)
		}
	}

}
