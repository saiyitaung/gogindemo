package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/gindemo/controller"
	"github.com/gindemo/service"
)

func SetViewsRoutes(r *gin.Engine, dbCon *sql.DB) {
	myCtl := controller.New(service.NewCategoryService(dbCon), service.NewProductService(dbCon), service.NewOrderService(dbCon))
	r.HandleMethodNotAllowed = true

	r.NoRoute(func(c *gin.Context) {
		c.HTML(404, "404.ghtml", nil)
	})

	r.GET("", myCtl.Home)
	r.GET("/:categoryName", myCtl.FindByCategory)
	r.GET("/cart", myCtl.ViewCart)
	r.POST("/addtocart", myCtl.AddToCart)
	r.POST("/removecartitem", myCtl.RemoveCartitem)
	r.POST("/updatcartitem", myCtl.UpdateCartItem)
	r.GET("/checkout", myCtl.CheckOutForm)
	r.POST("/checkout", myCtl.CheckOut)
	r.GET("/login", myCtl.LoginForm)
	r.POST("/login", myCtl.Login)

}
