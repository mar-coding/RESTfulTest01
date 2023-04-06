package main_test

import (
	"log"
	"os"
	"testing"

	marLib "github.com/mar-coding/RESTfulTest01"
)

var app marLib.App

func TestMain(m *testing.M) {
	dsn := marLib.LoadDB()
	app.Initialize("mysql", dsn)
	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func ensureTableExists() {
	if _, err := app.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func createTable() {
	// create my Movie table
	if _, err := app.DB.Exec(createPRKeyQuery); err != nil {
		log.Fatal(err)
	}
	if _, err := app.DB.Exec(addAutoIncrementQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	app.DB.Exec("DELETE FROM Movie")
	app.DB.Exec("ALTER TABLE Movie AUTO_INCREMENT = 1")
}

const tableCreationQuery = `CREATE TABLE Movie (
	movie_id int NOT NULL,
	movie_name varchar(200) NOT NULL,
	movie_year varchar(4) NOT NULL,
	movie_genre varchar(200) NOT NULL,
	movie_duration varchar(50) NOT NULL,
	movie_origin varchar(200) NOT NULL,
	movie_director varchar(200) NOT NULL,
	movie_rating float NOT NULL,
	movie_rating_count bigint NOT NULL,
	movie_link varchar(200) NOT NULL
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;`

const createPRKeyQuery = `ALTER TABLE Movie
ADD PRIMARY KEY (movie_id);`

const addAutoIncrementQuery = `ALTER TABLE Movie
MODIFY movie_id int NOT NULL AUTO_INCREMENT`
