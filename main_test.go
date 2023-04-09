package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	marLib "github.com/mar-coding/RESTfulTest01"
)

var app marLib.App

func TestUpdateMovie(t *testing.T) {
	// This test, first creates new random movie
	// then changes the detail of the movie and at last
	// check if it changes successfully or not

	clearTable()
	addMovies(1)

	req, _ := http.NewRequest("GET", "/movie/1", nil)
	res := executeRequest(req)
	var originalMovie map[string]interface{}
	json.Unmarshal(res.Body.Bytes(), &originalMovie)

	rawData := `{"id":1,"name":"The Updated Version","year":"9999","genre":"Test1 | Test2","duration":"3h 55min","origin":"USA","director":"Alice Bob","rate":1,"rate_count":1392322,"link":"https://www.test.com/title/xxx"}`
	var jsonStr = []byte(rawData)
	req, _ = http.NewRequest("PUT", "/movie/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	res = executeRequest(req)

	//check status code with 200
	checkResponseCode(t, http.StatusOK, res.Code)

	var m map[string]interface{}
	json.Unmarshal(res.Body.Bytes(), &m)

	if m["id"] != originalMovie["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalMovie["id"], m["id"])
	}

	if m["name"] == originalMovie["name"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalMovie["name"], m["name"], m["name"])
	}

	if m["price"] == originalMovie["price"] {
		t.Errorf("Expected the price to change from '%v' to '%v'. Got '%v'", originalMovie["price"], m["price"], m["price"])
	}

}
func TestEmptyTable(t *testing.T) {
	clearTable()
	req, _ := http.NewRequest("GET", "/movies", nil)
	res := executeRequest(req)

	//check status code with 200
	checkResponseCode(t, http.StatusOK, res.Code)
	if body := res.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}

}

func TestGetMovie(t *testing.T) {
	// This test first creates a new movie
	// then fetch it and checks the status code
	clearTable()
	addMovies(1)

	req, _ := http.NewRequest("GET", "/movie/1", nil)
	response := executeRequest(req)
	//check status code with 200
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestCreateMovie(t *testing.T) {
	// This test tries to create a test movie
	// and checks if the detail response is same as
	// the data of test movie
	clearTable()
	rawData := `{"id":1,"name":"The Test","year":"9999","genre":"Test1 | Test2","duration":"3h 55min","origin":"USA","director":"Alice Bob","rate":10,"rate_count":1392322,"link":"https://www.test.com/title/xxx"}`
	var jsonStr = []byte(rawData)
	req, _ := http.NewRequest("POST", "/movie", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	res := executeRequest(req)
	//check status code with 201
	checkResponseCode(t, http.StatusCreated, res.Code)

	var m map[string]interface{}
	json.Unmarshal(res.Body.Bytes(), &m)

	if m["name"] != "The Test" {
		t.Errorf("Expected movie name to be 'The Test'. Got '%v'", m["name"])
	}

	if m["rate"] != 10.0 {
		t.Errorf("Expected movie's rate to be '10'. Got '%v'", m["price"])
	}

	if m["id"] != 1.0 {
		t.Errorf("Expected movie ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetAbsentMovie(t *testing.T) {
	// This test tries to access a non-existent movie
	// check status code to be 404 & response body be
	// Movie not found.
	clearTable()
	url := fmt.Sprintf("/movie/%s", "10")
	req, _ := http.NewRequest("GET", url, nil)
	res := executeRequest(req)

	//check status code with 201
	checkResponseCode(t, http.StatusOK, res.Code)

	var m map[string]string
	json.Unmarshal(res.Body.Bytes(), &m)
	if m["error"] != "Movie not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Movie not found'. Got '%s'", m["error"])
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

func addMovies(num int) {
	// This function creates fake movie
	// the amount of movies that create is related on 'num' arg
	if num < 1 {
		num = 1
	}
	q := ""
	for i := 0; i < num; i++ {
		q = fmt.Sprintf("INSERT INTO Movie(movie_id, movie_name, movie_year, movie_genre, movie_duration, movie_origin,movie_director, movie_rating, movie_rating_count, movie_link) VALUES(1,%s,9999,Test1 | Test2,3h 55min,USA,Alice Bob,%d,1392322,www.test.com/title/xxx)", "Test Movie "+strconv.Itoa(i), rand.Intn(10-0)+0)
		app.DB.Exec(q)
	}
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
