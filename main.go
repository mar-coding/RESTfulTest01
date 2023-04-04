package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var counter int
var mutex = &sync.Mutex{}
var ENV_DB_USER string
var ENV_DB_PASSWORD string
var ENV_DB_IP string
var ENV_DB_PORT string
var ENV_DB_NAME string
var ENV_DB string

// create Data Source Name for connecting to DB
var DSN string

// movie represents data about a movie.
type Movie struct {
	ID    int     `json:"id"`
	Title string  `json:"title"`
	Year  string  `json:"year"`
	Genre string  `json:"genre"`
	Rate  float32 `json:"rate"`
}

// movies slice to seed movie data for testing API.
// var movies = []Movie{
// {ID: "0", Title: "The Shawshank Redemption", Year: "1994", Genre: "Drama", Rate: 9.3},
// {ID: "1", Title: "The Godfather", Year: "1972", Genre: "Crime, Drama", Rate: 9.2},
// {ID: "2", Title: "The Dark Knight", Year: "2008", Genre: "Action, Crime, Drama", Rate: 9.0},
// {ID: "3", Title: "The Godfather Part II", Year: "1974", Genre: "Crime, Drama", Rate: 9.0},
// {ID: "4", Title: "Schindler's List", Year: "1993", Genre: "Biography, Drama, History", Rate: 9.0},
// }
var movies = []Movie{}

//{"ID": "5", "Title": "TestTitle", "Year": "TestYear", "Genre": "Test, Test, Test", "Rate": 10}

// ---------------- REST section ----------------

// get all movies in json format
func getMovies(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(movies)
}

func addMovie(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var newMovie Movie
	json.Unmarshal(reqBody, &newMovie)
	movies = append(movies, newMovie)
	json.NewEncoder(w).Encode(newMovie)
}

func getMovieByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, _ := strconv.Atoi(vars["id"])
	for _, movie := range movies {
		if movie.ID == key {
			json.NewEncoder(w).Encode(movie)
		}
	}
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Extract the `id` of the movie we wish to delete
	id, _ := strconv.Atoi(vars["id"])
	// Loop through all our movies
	for index, movie := range movies {
		// if our id path parameter matches one of our
		// movies
		if movie.ID == id {
			// updates our movies array to remove the
			// movies
			movies = append(movies[:index], movies[index+1:]...)
		}
	}
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	//TODO: create update function to perform the updating
	// tasks
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	fmt.Println("Endpoint Hit: homePage")
}

func aminPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi")
	fmt.Println("Endpoint Hit: aminPage")
}

func incrementCounterPage(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	counter++
	fmt.Fprintf(w, strconv.Itoa(counter))
	mutex.Unlock()
	fmt.Println("Endpoint Hit: incPage")
}

func handleRequests() {

	//API version variable
	api_ver := "/api/v1"

	// - /api/v1/post - HTTP GET request - All Posts
	// - /api/v1/post/:id - HTTP GET request - Single post
	// - /api/v1/post/:id - HTTP POST request - Publish a post
	// - /api/v1/post/:id - HTTP PATCH request - update an existing post
	// - /api/v1/post/:id - HTTP DELETE request - deletes a post

	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/amin", aminPage)
	myRouter.HandleFunc("/inc", incrementCounterPage)
	myRouter.HandleFunc(fmt.Sprintf("%s/movie", api_ver), getMovies)
	// it only calls the addMovie when http method is a POST
	// NOTE: Ordering is important here! This has to be defined before
	// the other `/movie` endpoint.
	myRouter.HandleFunc(fmt.Sprintf("%s/movie", api_ver), addMovie).Methods("POST")
	myRouter.HandleFunc(fmt.Sprintf("%s/movie/{id}", api_ver), deleteMovie).Methods("DELETE")
	myRouter.HandleFunc(fmt.Sprintf("%s/movie/{id}", api_ver), updateMovie).Methods("PUT")
	myRouter.HandleFunc(fmt.Sprintf("%s/movie/{id}", api_ver), getMovieByID)
	// http.Handle("/", http.FileServer(http.Dir("./static")))
	myRouter.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	// http.HandleFunc("/", indexPage)

	log.Fatal(http.ListenAndServe(":12345", myRouter))

	// for https running
	// log.Fatal(http.ListenAndServeTLS(":443", "server.crt", "server.key", nil))

}

// ---------------- DB section ----------------
func load_db_env() {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatalf("Some error occured at loading .ENV file. Err: %s", err)
	}
	ENV_DB_USER = os.Getenv("DB_USER")
	ENV_DB_PASSWORD = os.Getenv("DB_PASSWORD")
	ENV_DB_IP = os.Getenv("DB_IP")
	ENV_DB_PORT = os.Getenv("DB_PORT")
	ENV_DB_NAME = os.Getenv("DB_NAME")
	ENV_DB = os.Getenv("DB")

	DSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", ENV_DB_USER, ENV_DB_PASSWORD, ENV_DB_IP, ENV_DB_PORT, ENV_DB_NAME)
}
func db_connect() (db *sql.DB) {
	db, err := sql.Open(ENV_DB, DSN)
	if err != nil {
		panic(err.Error())
		// log.Fatalf("impossible to create the connection: %s", err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db
}
func db_insert(movie_name string, movie_year string, movie_genre string, movie_rate string) {
	fmt.Println("** Insert **")
	db := db_connect()
	defer db.Close()
	// perform insert
	q := fmt.Sprintf("INSERT INTO Movie(movie_id, movie_name, movie_year, movie_genre, movie_rate) VALUES (?, ?, ?, ?, ?)")
	insertResult, err := db.ExecContext(context.Background(), q, nil, movie_name, movie_year, movie_genre, movie_rate)
	if err != nil {
		log.Fatalf("impossible insert movie: %s", err)
	}
	id, err := insertResult.LastInsertId()
	if err != nil {
		log.Fatalf("impossible to retrieve last inserted id: %s", err)
	}
	log.Printf("inserted id: %d", id)

}
func db_retrieve_all() []Movie {
	fmt.Println("** Retrieve all data from database **")
	db := db_connect()
	defer db.Close()
	// perform retrieve
	q := fmt.Sprintf("SELECT * FROM Movie ORDER BY movie_id ASC")
	selDB, err := db.Query(q)
	if err != nil {
		panic(err.Error())
	}
	output := []Movie{}
	temp := Movie{}
	for selDB.Next() {
		var id int
		var name, year, genre string
		var rate float32
		err = selDB.Scan(&id, &name, &year, &genre, &rate)
		if err != nil {
			panic(err.Error())
		}
		temp.ID = id
		temp.Title = name
		temp.Year = year
		temp.Genre = genre
		temp.Rate = rate
		output = append(output, temp)
	}
	return output
}
func db_retrieve_by_id(id int) Movie {
	fmt.Println("** Retrieve data based on ID **")
	db := db_connect()
	defer db.Close()

	selDB, err := db.Query("SELECT * FROM Movie WHERE movie_id=?", id)
	if err != nil {
		panic(err.Error())
	}
	temp := Movie{}
	for selDB.Next() {
		var id int
		var name, year, genre string
		var rate float32
		err = selDB.Scan(&id, &name, &year, &genre, &rate)
		if err != nil {
			panic(err.Error())
		}
		temp.ID = id
		temp.Title = name
		temp.Year = year
		temp.Genre = genre
		temp.Rate = rate
	}
	return temp
}
func db_update(movie_id int, movie_name string, movie_year string, movie_genre string, movie_rate string) {
	fmt.Println("** Update data based on ID **")
	db := db_connect()
	defer db.Close()
	q, err := db.Prepare("UPDATE Movie SET movie_name=?, movie_year=?, movie_genre=?, movie_rate=? WHERE movie_id=?")
	if err != nil {
		panic(err.Error())
	}
	q.Exec(movie_name, movie_year, movie_genre, movie_rate, movie_id)
}
func db_delete(movie_id int) {
	fmt.Println("** Delete data based on ID **")
	db := db_connect()
	defer db.Close()
	q, err := db.Prepare("DELETE FROM Movie WHERE movie_id=?")
	if err != nil {
		panic(err.Error())
	}
	q.Exec(movie_id)
}

func main() {
	fmt.Println("Rest API v1.0 - with Mux Routers")
	load_db_env()
	// db_insert("The Dark Knight", "2008", "Action, Crime, Drama", "9.0")
	// db_insert("The Godfather Part II", "1974", "Crime, Drama", "9.0")
	// db_insert("Schindler's List", "1993", "Biography, Drama, History", "9.0")
	// db_update(10, "The Godfather Part II", "1974", "Crime, Drama", "9.0")
	db_delete(11)
	t := db_retrieve_all()
	// t := db_retrieve_by_id(5)
	fmt.Println(t)
	// handleRequests()
}
