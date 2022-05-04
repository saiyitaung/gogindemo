package middlewares

import "github.com/dgrijalva/jwt-go"

type MyCustomClaims struct {
	jwt.StandardClaims
	Data string `json:"data"`
}
