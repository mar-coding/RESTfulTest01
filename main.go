package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
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
// movie_id, movie_name, movie_year, movie_genre, movie_duration, movie_origin, movie_director, movie_rate, movie_rate_count, movie_link
type Movie struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Year       string  `json:"year"`
	Genre      string  `json:"genre"`
	Duration   string  `json:"duration"`
	Origin     string  `json:"origin"`
	Director   string  `json:"director"`
	Rate       float32 `json:"rate"`
	Rate_count uint64  `json:"rate_count"`
	Link       string  `json:"link"`
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
func db_insert(m Movie) {
	fmt.Println("** Insert **")
	db := db_connect()
	defer db.Close()
	// perform insert
	q := fmt.Sprintf("INSERT INTO Movie(movie_id, movie_name, movie_year, movie_genre, movie_duration, movie_origin, movie_director, movie_rating, movie_rating_count, movie_link) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	insertResult, err := db.ExecContext(context.Background(), q, m.ID, m.Name, m.Year, m.Genre, m.Duration, m.Origin, m.Director, m.Rate, m.Rate_count, m.Link)
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
		var name, year, genre, duration, origin, director, link string
		var rate float32
		var rate_count uint64
		err = selDB.Scan(&id, &name, &year, &genre, &duration, &origin, &director, &rate, &rate_count, &link)
		if err != nil {
			panic(err.Error())
		}
		temp.ID = id
		temp.Name = name
		temp.Year = year
		temp.Genre = genre
		temp.Duration = duration
		temp.Origin = origin
		temp.Director = director
		temp.Rate = rate
		temp.Rate_count = rate_count
		temp.Link = link
		output = append(output, temp)
	}
	return output
}
func db_retrieve_by_id(id int) (Movie, error) {
	fmt.Println("** Retrieve data based on ID **")
	db := db_connect()
	defer db.Close()

	exist := 0
	q, err := db.Query("SELECT EXISTS(SELECT * from Movie WHERE movie_id=?) AS temp; ", id)
	for q.Next() {
		err = q.Scan(&exist)
	}

	if err != nil {
		panic(err.Error())
	}
	temp := Movie{}
	if exist == 1 {
		selDB, err := db.Query("SELECT * FROM Movie WHERE movie_id=?", id)
		// fmt.Println(err)
		if err != nil {
			panic(err.Error())
		}
		for selDB.Next() {
			var id int
			var name, year, genre, duration, origin, director, link string
			var rate float32
			var rate_count uint64
			err = selDB.Scan(&id, &name, &year, &genre, &duration, &origin, &director, &rate, &rate_count, &link)
			if err != nil {
				panic(err.Error())
			}
			temp.ID = id
			temp.Name = name
			temp.Year = year
			temp.Genre = genre
			temp.Duration = duration
			temp.Origin = origin
			temp.Director = director
			temp.Rate = rate
			temp.Rate_count = rate_count
			temp.Link = link
		}
		return temp, nil
	} else {
		return temp, errors.New("id not found.")
	}
}
func db_update(m Movie) {
	fmt.Println("** Update data based on ID **")
	db := db_connect()
	defer db.Close()
	q, err := db.Prepare("UPDATE Movie SET movie_name=?, movie_year=?, movie_genre=?, movie_duration=?, movie_origin=?, movie_director=?, movie_rating=?, movie_rating_count=?, movie_link=? WHERE movie_id=?")
	if err != nil {
		panic(err.Error())
	}
	q.Exec(m.Name, m.Year, m.Genre, m.Duration, m.Origin, m.Director, m.Rate, m.Rate_count, m.Link, m.ID)
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
	// m := Movie{1, "The Shawshank Redemption", "1994", "Drama", "2h 22min", "USA", "Frank Darabont", 9.3, 2030817, "https://www.imdb.com/title/tt0111161"}
	// db_insert(m)
	// db_update(m)
	// db_delete(11)
	// t, err := db_retrieve_by_id(1)
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(t)
	// }
	// handleRequests()
	// t := db_retrieve_all()
	// fmt.Println(t)
}
