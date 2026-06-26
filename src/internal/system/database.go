package system

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func OpenDatabase() (*sql.DB, error) {
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	name := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")

	db, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, name),
	)

	return db, err
}
