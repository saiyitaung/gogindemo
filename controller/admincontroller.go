package controller

import (
	"crypto/md5"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gindemo/entities"
	"github.com/gindemo/service"
	uuid "github.com/satori/go.uuid"
)

type AdminController interface {
	AdminHome(ctx *gin.Context)
	AdminProducts(ctx *gin.Context)
	AdminProductsPost(ctx *gin.Context)
	AdminProductsNewForm(ctx *gin.Context)
	UploadPic(ctx *gin.Context)
	Customers(ctx *gin.Context)
	GetOrders(ctx *gin.Context)
	GetOrderDetail(ctx *gin.Context)
	AdminCategoryNewForm(ctx *gin.Context)
	AdminCategory(ctx *gin.Context)
	CreateCategory(ctx *gin.Context)
	Logout(ctx *gin.Context)
	Error(c *gin.Context)
}
type adminCtrlImpl struct {
	categoryService service.CategoryDaoService
	orderService    service.OrderDaoService
	customerService service.CustomerDaoService
	rpService       service.ReportDaoService
	pService        service.ProductDaoService
}

func NewAdminCtrl(cats service.CategoryDaoService, os service.OrderDaoService, cus service.CustomerDaoService, rps service.ReportDaoService, ps service.ProductDaoService) AdminController {
	return &adminCtrlImpl{
		categoryService: cats,
		orderService:    os,
		customerService: cus,
		rpService:       rps,
		pService:        ps,
	}
}

func (admin *adminCtrlImpl) AdminHome(ctx *gin.Context) {
	m, err := admin.rpService.GetAllSales()
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	topSales, err := admin.rpService.AllSalesByCategories()
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	wSales, err := admin.rpService.SalesInWeek(time.Now())
	if err != nil {
		fmt.Println(err)
	}
	data := struct {
		ProductSales  map[string]int
		CategorySales map[string]float64
		WeeklySales   map[string]float64
	}{
		ProductSales:  m,
		CategorySales: topSales,
		WeeklySales:   wSales,
	}

	ctx.HTML(200, ADMIN_HOME, Payload{Title: "Admin", Data: data})
}
func (admin *adminCtrlImpl) AdminProducts(ctx *gin.Context) {
	products, err := admin.pService.FindAll()
	if err != nil {
		fmt.Println(err)
	}
	ctx.HTML(200, ADMIN_PRODUCTS, Payload{Title: "Products", Data: products})
}
func (admin *adminCtrlImpl) Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Save()
	ctx.Redirect(http.StatusSeeOther, "/login")
}
func (admin *adminCtrlImpl) AdminProductsPost(ctx *gin.Context) {
	cid := ctx.PostForm("cid")
	pname := ctx.PostForm("pname")
	pprice := ctx.PostForm("pprice")
	pqty := ctx.PostForm("pqty")
	pdesc := ctx.PostForm("pdesc")
	imgs := ctx.PostFormArray("imgpaths")
	fmt.Println(imgs)
	h := md5.New()
	h.Write([]byte(time.Now().String()))
	var p = entities.Product{ID: uuid.NewV4().String(), Name: pname, CategoryId: cid, Description: pdesc, Images: imgs, CoverPic: imgs[0], Created: time.Now(), LastUpdate: time.Now()}
	price, err := strconv.ParseFloat(pprice, 64)
	if err != nil {
		ctx.Set("errMsg", err.Error())
		ctx.Redirect(http.StatusSeeOther, ERR_ROUTE)
		return
	}
	qty, err := strconv.Atoi(pqty)
	if err != nil {
		ctx.Set("errMsg", err.Error())
		ctx.Redirect(http.StatusSeeOther, ERR_ROUTE)
		return
	}
	p.Price = price
	p.Qty = qty
	// mc.products = append(mc.products, p)
	_, err = admin.pService.Create(p)
	if err != nil {
		ctx.Set("errMsg", err.Error())
		ctx.Redirect(http.StatusSeeOther, ERR_ROUTE)
		return
	}
	ctx.Redirect(303, ADMIN_PRODUCTS_ROUTE)
}

func (admin *adminCtrlImpl) AdminProductsNewForm(ctx *gin.Context) {
	categories, _ := admin.categoryService.FindAll()
	ctx.HTML(200, ADMIN_PRODUCTNEW, struct {
		Title      string
		Categories []entities.Category
	}{Title: "New Product", Categories: categories})
}

func (admin *adminCtrlImpl) UploadPic(ctx *gin.Context) {

	f, fh, err := ctx.Request.FormFile("mypic")
	if err != nil {
		ctx.JSON(500, gin.H{"msg": "Error"})
		return
	}
	fmt.Println("", fh.Filename)
	fmt.Println("", fh.Size)
	h := md5.New()
	h.Write([]byte("product"))
	ext := strings.Split(fh.Filename, ".")[1]
	fn := fmt.Sprintf("%x", h.Sum([]byte(fh.Filename))) + "." + ext
	upf, err := os.Create(filepath.Join(UPLOAD_DIR, fn))
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"msg": "Error"})
		return
	}
	io.Copy(upf, f)
	ctx.JSON(http.StatusCreated, fn)
}
func (admin *adminCtrlImpl) Customers(ctx *gin.Context) {
	pageStr := ctx.Query("page")
	isFirstPage := false
	isLastPage := false
	currentPage, err := strconv.Atoi(pageStr)
	if err != nil {
		fmt.Println(err)
		currentPage = 1
	}
	total, err := admin.customerService.TotalRecords()
	if err != nil {
		ctx.Redirect(http.StatusSeeOther, "/admin/error")
		return
	}
	limit := 15
	numbersOfPages := math.Ceil((float64(total)) / float64(limit))
	if currentPage > int(numbersOfPages) {
		admin.Error(ctx)
		return
	}
	var m = make(map[string]interface{})
	if currentPage > 1 {
		m["Prev"] = currentPage - 1
	} else {
		m["Prev"] = currentPage
	}
	if currentPage < int(numbersOfPages) {
		m["Next"] = currentPage + 1
	} else {
		m["Next"] = currentPage
	}
	if currentPage == 1 {
		isFirstPage = true
	}
	if currentPage == int(numbersOfPages) {
		isLastPage = true
	}

	if currentPage > 0 && currentPage <= int(numbersOfPages) {
		customers, err := admin.customerService.FindAll((currentPage-1)*limit, limit)
		if err != nil {
			fmt.Println(err)
			ctx.Abort()
			return
		}

		m["Customers"] = customers
		m["IsFirst"] = isFirstPage
		m["IsLast"] = isLastPage
		ctx.HTML(200, ADMIN_CUSTOMER, Payload{Title: "Customers", Data: m})
		return
	}
	customers, err := admin.customerService.FindAll(0, limit)
	m["Customers"] = customers
	m["IsFirst"] = isFirstPage
	m["IsLast"] = isLastPage
	if err != nil {
		fmt.Println(err)
		ctx.Abort()
		return
	}
	ctx.HTML(200, ADMIN_CUSTOMER, Payload{Title: "Customers", Data: m})
}
func (admin *adminCtrlImpl) GetOrders(ctx *gin.Context) {
	pageStr := ctx.Query("page")
	currentPage, err := strconv.Atoi(pageStr)
	isFirstPage := false
	isLastPage := false
	if err != nil {
		fmt.Println(err)
		currentPage = 1
	}
	total, err := admin.orderService.GetTotalCount()
	if err != nil {
		ctx.Redirect(http.StatusSeeOther, "/admin/error")
		return
	}
	limit := 15
	numbersOfPages := math.Ceil((float64(total)) / float64(limit))
	fmt.Println("nums of page ", numbersOfPages)
	fmt.Println("current page", currentPage)
	var m = make(map[string]interface{})
	if currentPage > int(numbersOfPages) {
		admin.Error(ctx)
		return
	}
	if currentPage > 1 {
		m["Prev"] = currentPage - 1
	} else {
		m["Prev"] = currentPage
	}
	if currentPage < int(numbersOfPages) {
		m["Next"] = currentPage + 1
	} else {
		m["Next"] = currentPage
	}
	if currentPage == 1 {
		isFirstPage = true
	}
	if currentPage == int(numbersOfPages) {
		isLastPage = true
	}
	if currentPage > 0 && currentPage <= int(numbersOfPages) {
		orders, err := admin.orderService.FindAll((currentPage-1)*limit, limit)
		if err != nil {
			fmt.Println(err)
			ctx.Abort()
			return
		}

		m["Orders"] = orders
		m["IsFirst"] = isFirstPage
		m["IsLast"] = isLastPage
		ctx.HTML(200, ADMIN_ORDERS, Payload{Title: "Orders", Data: m})
		return
	}

	orders, err := admin.orderService.FindAll(0, limit)
	if err != nil {
		fmt.Println(err)
		ctx.Abort()
		return
	}
	m["Orders"] = orders
	m["IsFirst"] = isFirstPage
	m["IsLast"] = isLastPage
	ctx.HTML(200, ADMIN_ORDERS, Payload{Title: "Orders", Data: m})
}
func (admin *adminCtrlImpl) GetOrderDetail(ctx *gin.Context) {
	orderId := ctx.Params.ByName("id")
	orderDetail, err := admin.orderService.OrderDetail(orderId)
	if err != nil {
		fmt.Println(err)
		ctx.Abort()
		return
	}
	ctx.HTML(200, ADMIN_ORDERS_DETAIL, Payload{Title: "Order Detail", Data: orderDetail})
}
func (admin *adminCtrlImpl) AdminCategoryNewForm(ctx *gin.Context) {
	ctx.HTML(200, ADMIN_CATEGORYNEW, struct{ Title string }{"New Category"})
}
func (admin *adminCtrlImpl) AdminCategory(ctx *gin.Context) {
	categories, err := admin.categoryService.FindAll()
	if err != nil {
		fmt.Println(err)
	}
	ctx.HTML(200, ADMIN_CATEGORY, struct {
		Title string
		Data  interface{}
	}{Title: "Categories", Data: categories})
}

func (admin *adminCtrlImpl) CreateCategory(ctx *gin.Context) {
	cname := ctx.PostForm("cname")
	session := sessions.Default(ctx)
	h := md5.New()
	if strings.TrimLeft(cname, " ") == "" {
		session.Set("errMsg", "Category Name can not be empty")
		session.Save()
		ctx.Redirect(http.StatusSeeOther, ERR_ROUTE)
		return
	}
	f, fh, err := ctx.Request.FormFile("categorypic")
	if err != nil {
		session.Set("errMsg", "Category Image can not be empty")
		session.Save()
		ctx.Redirect(http.StatusSeeOther, ERR_ROUTE)
		return
	}
	h.Write([]byte(time.Now().String()))
	id := fmt.Sprintf("%x", h.Sum([]byte(cname)))
	var newCategory = entities.Category{ID: uuid.NewV4().String(), Name: cname}
	ext := strings.Split(fh.Filename, ".")[1]
	dest, err := os.Create(filepath.Join(UPLOAD_DIR, id+"."+ext))
	if err != nil {
		session.Set("errMsg", err.Error())
		session.Save()
		ctx.Redirect(http.StatusSeeOther, ERR_ROUTE)
		return
	}
	_, err = io.Copy(dest, f)
	if err != nil {
		session.Set("errMsg", err.Error())
		session.Save()
		ctx.Redirect(http.StatusSeeOther, ERR_ROUTE)
		return
	}
	newCategory.Image = "http://localhost:9110/resources/imgs/" + id + "." + ext
	fmt.Println(newCategory)
	_, err = admin.categoryService.Create(newCategory)
	if err != nil {
		session.Set("errMsg", err.Error())
		session.Save()
		ctx.Redirect(http.StatusSeeOther, ERR_ROUTE)
		return
	}
	ctx.Redirect(http.StatusSeeOther, ADMIN_CATEGORY_ROUTE)
}
func (admin *adminCtrlImpl) Error(c *gin.Context) {
	s := sessions.Default(c)
	if msg, ok := s.Get("errMsg").(string); ok {
		c.HTML(200, "error.ghtml", ErrData{ErrMsg: msg})
		return
	}
	c.HTML(200, "error.ghtml", ErrData{"Unknow Error"})
}

type ErrData struct {
	ErrMsg interface{}
}
