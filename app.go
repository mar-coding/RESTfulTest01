package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(db_name string, dsn string) {
	a.DB = dbConnect(db_name, dsn)
	a.Router = mux.NewRouter()
}

func (a *App) Run(addr string) {}

func dbConnect(db_name string, dsn string) (db *sql.DB) {
	db, err := sql.Open(db_name, dsn)
	if err != nil {
		log.Fatal(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return
}
func LoadDB() (dsn string) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured at loading .ENV file. Err: %s", err)
	}
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbIp := os.Getenv("DB_IP")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbIp, dbPort, dbName)
	return
}
