package controller

import (
	"fmt"
	"net/http"
	"time"

	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gindemo/entities"
	mymiddleware "github.com/gindemo/middlewares"
	"github.com/gindemo/service"
	myutils "github.com/gindemo/utils"
	uuid "github.com/satori/go.uuid"
)

type MyController interface {
	Home(ctx *gin.Context)
	ViewCart(ctx *gin.Context)
	CheckOutForm(ctx *gin.Context)
	CheckOut(ctx *gin.Context)
	UpdateCartItem(ctx *gin.Context)
	RemoveCartitem(ctx *gin.Context)
	AddToCart(ctx *gin.Context)
	FindByCategory(ctx *gin.Context)

	LoginForm(ctx *gin.Context)
	Login(ctx *gin.Context)
}
type myController struct {
	orderService    service.OrderDaoService
	categoryService service.CategoryDaoService
	productService  service.ProductDaoService
}

func New(cs service.CategoryDaoService, ps service.ProductDaoService, orderS service.OrderDaoService) MyController {
	return &myController{categoryService: cs, productService: ps, orderService: orderS}
}
func (mc *myController) Home(ctx *gin.Context) {

	cagetories, err := mc.categoryService.FindAll()
	if err != nil {
		fmt.Println("Err ", err)
	}
	products, err := mc.productService.FindAll()
	if err != nil {
		fmt.Println(err)
	}
	Data := struct {
		Items      int
		Categories []entities.Category
		Products   []entities.Product
	}{
		Items:      0,
		Categories: cagetories,
		Products:   products,
	}
	cart := sessions.Default(ctx)
	if m, ok := cart.Get("cartItems").(map[string]int); ok {
		fmt.Println(m)
		Data.Items = len(m)
	}
	payload := Payload{Title: "Gin Web Framework", Data: Data}
	ctx.HTML(200, INDEX, payload)
}
func (admin *myController) LoginForm(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.ghtml", Payload{Title: "Login"})
}
func (admin *myController) Login(ctx *gin.Context) {
	username := ctx.PostForm("username")
	pass := ctx.PostForm("password")
	if username != "user" || pass != "admin" {
		ctx.Redirect(http.StatusSeeOther, "/login")
		return
	}
	claims := mymiddleware.MyCustomClaims{StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Add(time.Minute * 10).Unix(), Issuer: "test"}, Data: "Data"}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(myutils.GetSecretKey())
	if err != nil {
		fmt.Println(err)
		ctx.JSON(200, gin.H{"err": "Sign token err"})
		return
	}
	session := sessions.Default(ctx)
	fmt.Println(tokenString)
	session.Set("sessionId", tokenString)
	session.Save()
	ctx.Redirect(http.StatusSeeOther, "/admin")
}
func (mc *myController) ViewCart(ctx *gin.Context) {
	cartSession := sessions.Default(ctx)
	if m, ok := cartSession.Get("cartItems").(map[string]int); ok {
		var cart = entities.Cart{Items: []entities.CartItem{}, Total: 0.0}
		for k, v := range m {
			p, err := mc.productService.FindById(k)
			if err != nil {
				fmt.Println("Get Product Error :", err)
				continue
			}
			var ci = entities.CartItem{Product: p, Count: v}
			ci.CalcPrice()
			cart.Items = append(cart.Items, ci)

		}
		cart.CalcTotal()
		ctx.HTML(200, "cart.ghtml", Payload{Title: "Cart", Data: cart})
		return
	}
	ctx.HTML(200, "cart.ghtml", Payload{Title: "Cart", Data: entities.Cart{Items: []entities.CartItem{}}})
}
func (mc *myController) CheckOutForm(ctx *gin.Context) {
	cartSession := sessions.Default(ctx)
	if m, ok := cartSession.Get("cartItems").(map[string]int); ok {
		var cart = entities.Cart{Items: []entities.CartItem{}, Total: 0.0}
		for k, v := range m {
			p, err := mc.productService.FindById(k)
			if err != nil {
				fmt.Println("Get Product Error :", err)
				continue
			}
			var ci = entities.CartItem{Product: p, Count: v}
			ci.CalcPrice()
			cart.Items = append(cart.Items, ci)

		}
		cart.CalcTotal()
		if len(m) > 0 {
			d := struct {
				Count int
				Cart  entities.Cart
			}{
				Count: len(cart.Items),
				Cart:  cart,
			}
			ctx.HTML(200, "checkoutform.ghtml", Payload{Title: "Checkout Form", Data: d})
			return
		}
	}
	ctx.Redirect(http.StatusSeeOther, "/")
}
func (mc *myController) CheckOut(ctx *gin.Context) {

	name := ctx.PostForm("name")
	email := ctx.PostForm("email")
	phone := ctx.PostForm("phone")
	addr := ctx.PostForm("address")
	var c = entities.Customer{ID: uuid.NewV4().String(), Name: name, Email: email, Phone: phone, Address: addr}
	fmt.Println(c)
	cartSession := sessions.Default(ctx)
	if m, ok := cartSession.Get("cartItems").(map[string]int); ok {
		var cart = entities.Cart{Items: []entities.CartItem{}, Total: 0.0}
		for k, v := range m {
			p, err := mc.productService.FindById(k)
			if err != nil {
				fmt.Println("Get Product Error :", err)
				continue
			}
			var ci = entities.CartItem{Product: p, Count: v}
			ci.CalcPrice()
			cart.Items = append(cart.Items, ci)

		}
		cart.CalcTotal()
		var newOrder = entities.Orders{ID: uuid.NewV4().String(), Created: time.Now(), NumberOfProducts: len(cart.Items), Amount: cart.Total, Confirm: true}
		o_detail, err := mc.orderService.CreateCustomerOrder(c, newOrder, cart.Items)
		if err != nil {
			fmt.Println(err)
			ctx.Abort()
			return
		}
		fmt.Println(o_detail)

	}

	cartSession.Clear()
	cartSession.Save()
	ctx.Redirect(http.StatusSeeOther, "/")
}
func (mc *myController) UpdateCartItem(ctx *gin.Context) {
	pId := ctx.PostForm("productId")
	count := ctx.PostForm("count")
	intCount, err := strconv.Atoi(count)
	if err != nil {
		fmt.Println(err)
	}
	cartSession := sessions.Default(ctx)
	if m, ok := cartSession.Get("cartItems").(map[string]int); ok {
		if v, ok := m[pId]; ok {
			v += intCount
			m[pId] = v
		}
		cartSession.Set("cartItems", m)
		cartSession.Save()
	}
	ctx.Redirect(http.StatusSeeOther, "/cart")
}
func (mc *myController) RemoveCartitem(ctx *gin.Context) {
	pId := ctx.PostForm("productId")
	count := ctx.PostForm("count")
	intCount, err := strconv.Atoi(count)
	if err != nil {
		fmt.Println(err)
	}
	cartSession := sessions.Default(ctx)

	if m, ok := cartSession.Get("cartItems").(map[string]int); ok {
		if v, ok := m[pId]; ok {
			v -= intCount
			if v > 0 {
				m[pId] = v
			} else {
				delete(m, pId)
			}
		}
		cartSession.Set("cartItems", m)
		cartSession.Save()
	}
	ctx.Redirect(http.StatusSeeOther, "/cart")
}
func (mc *myController) AddToCart(ctx *gin.Context) {
	pId := ctx.PostForm("productId")
	itemCount := ctx.PostForm("count")
	fmt.Println(pId, itemCount)
	intCount, err := strconv.Atoi(itemCount)
	if err != nil {
		fmt.Println(err)
		ctx.Abort()
		return
	}

	cartSession := sessions.Default(ctx)
	if m, ok := cartSession.Get("cartItems").(map[string]int); !ok {
		cartSession.Set("cartItems", map[string]int{pId: intCount})
		cartSession.Options(sessions.Options{MaxAge: 60 * 30})
		err = cartSession.Save()
		if err != nil {
			fmt.Println(err)
		}

	} else {

		var v = m[pId]
		m[pId] = v + intCount
		cartSession.Set("cartItems", m)
		err = cartSession.Save()

		if err != nil {
			fmt.Println(err)
		}
	}
	// cartSession.Options(sessions.Options{})

	ctx.Redirect(http.StatusSeeOther, "/")
}
func (mc *myController) FindByCategory(ctx *gin.Context) {
	param := ctx.Params.ByName("categoryName")
	categories, err := mc.categoryService.FindAll()
	if err != nil {
		fmt.Println(err)
		ctx.Abort()
		return
	}
	for _, c := range categories {
		if c.Name == param {
			products, err := mc.productService.FindByCategoryID(c.ID)
			if err != nil {
				fmt.Println(err)
				ctx.Abort()
				return
			}
			Data := struct {
				Items      int
				Categories []entities.Category
				Products   []entities.Product
			}{
				Items:      0,
				Categories: categories,
				Products:   products,
			}
			cart := sessions.Default(ctx)
			if m, ok := cart.Get("cartItems").(map[string]int); ok {
				fmt.Println(m)
				Data.Items = len(m)
			}
			ctx.HTML(200, INDEX, Payload{Title: param, Data: Data})
			return
		}
	}
	ctx.Abort()
}

type Payload struct {
	Title string
	Data  interface{}
}
