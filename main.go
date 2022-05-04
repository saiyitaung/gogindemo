package main

import (
	"encoding/gob"
	"fmt"

	_ "github.com/lib/pq"

	"html/template"
	"log"
	"net/http"

	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/gindemo/controller"
	"github.com/gindemo/routes"

	myutils "github.com/gindemo/utils"
)

func main() {
	engine := gin.Default()

	store := cookie.NewStore(myutils.GetSecretKey())
	gob.Register(map[string]int{})
	var fm = template.FuncMap{
		"shortTxt": func(s string) string {
			if len(s) > 100 {
				return s[:100] + "..."
			}
			return s
		},
		"gIndex": func(i int) int {
			return i + 1
		},
		"dateFmt": func(t time.Time) string {
			return t.Format("Mon , 01-02-2006")
		},
		"floatFmt": func(p float64) string {
			return fmt.Sprintf("%.2f", p)
		},
	}
	engine.SetFuncMap(fm)

	engine.LoadHTMLGlob("./public/templates/*.ghtml")
	engine.StaticFS(controller.RESOURCES_ROUTE, http.Dir("./public/resources/"))
	engine.Use(sessions.Sessions("cart", store))

	db, err := myutils.GetDB()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		closeEr := db.Close()
		if closeEr != nil {
			log.Println("closeing db Err :", closeEr)
		}
	}()
	routes.SetViewsRoutes(engine, db)
	routes.SetAdminRoutes(engine.Group("/admin"), db)
	routes.SetApiRoutes(engine.Group("/api"), db)

	log.Fatal(engine.Run(":9110"))
}
