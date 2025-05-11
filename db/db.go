package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func InitDB() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	addr := os.Getenv("DB_ADDR")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_NAME")

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, addr, port, database)

	DB, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalln(err)
	}

}
