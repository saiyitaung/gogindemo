package utils

import "os"

func GetSecretKey() []byte {
	secret := os.Getenv("SECRET")
	if secret == "" {
		secret = "supersecret"
	}
	return []byte(secret)
}
