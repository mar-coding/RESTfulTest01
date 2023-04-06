package main_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	marLib "github.com/mar-coding/RESTfulTest01"
)

var app marLib.App

func TestEmptyTable(t *testing.T) {
	clearTable()
	req, _ := http.NewRequest("GET", "/movies", nil)
	res := executeRequest(req)

	checkResponseCode(t, http.StatusOK, res.Code)
	if body := res.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}

}

func TestMain(m *testing.M) {
	dsn := marLib.LoadDB()
	app.Initialize("mysql", dsn)
	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func executeRequest(req *http.Request) (rr *httptest.ResponseRecorder) {
	//executes the request and callbacks the response
	rr = httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)
	return
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
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
	// it will clear all data on Movie table
	app.DB.Exec("DELETE FROM Movie")
	app.DB.Exec("ALTER TABLE Movie AUTO_INCREMENT = 1")
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS Movie (
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
