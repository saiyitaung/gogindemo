package utils

import (
	"database/sql"
	"fmt"
	"os"
)

func initDB(user, pass string) (*sql.DB, error) {
	return sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@localhost/moviedb?sslmode=disable", user, pass))
}
func GetDB() (*sql.DB, error) {
	return initDB(os.Getenv("PSQLU"), os.Getenv("PSQLP"))
}
