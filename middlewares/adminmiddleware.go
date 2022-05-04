package middlewares

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	myutils "github.com/gindemo/utils"
)

func IsAuthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		tokenStr := session.Get("sessionId")
		if tokenStr == nil {
			fmt.Println("No Token :: ")
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}
		//fmt.Println(tokenStr.(string))
		token, err := jwt.ParseWithClaims(tokenStr.(string), &MyCustomClaims{}, func(t *jwt.Token) (interface{}, error) {
			return myutils.GetSecretKey(), nil
		})
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}
		if _, ok := token.Claims.(*MyCustomClaims); !ok && !token.Valid {
			fmt.Println("Invalid claims")
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}
		fmt.Println("middleware : ", c.Request.URL.Path)
	}
}
