package sqlconnect

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() (*sql.DB, error) {
	user := os.Getenv("DB_USER")
	passsword := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	dbport := os.Getenv("DB_PORT")
	host := os.Getenv("HOST")

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, passsword, host, dbport, dbname)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}
	return db, nil
}
