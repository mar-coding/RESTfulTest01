package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

var counter int
var mutex = &sync.Mutex{}

// movie represents data about a movie.
type Movie struct {
	ID    string  `json:"id"`
	Title string  `json:"title"`
	Year  string  `json:"year"`
	Genre string  `json:"genre"`
	Rate  float32 `json:"rate"`
}

// movies slice to seed movie data.
var movies = []Movie{
	{ID: "0", Title: "The Shawshank Redemption", Year: "1994", Genre: "Drama", Rate: 9.3},
	{ID: "1", Title: "The Godfather", Year: "1972", Genre: "Crime, Drama", Rate: 9.2},
	{ID: "2", Title: "The Dark Knight", Year: "2008", Genre: "Action, Crime, Drama", Rate: 9.0},
	{ID: "3", Title: "The Godfather Part II", Year: "1974", Genre: "Crime, Drama", Rate: 9.0},
	{ID: "4", Title: "Schindler's List", Year: "1993", Genre: "Biography, Drama, History", Rate: 9.0},
}

//{"ID": "5", "Title": "TestTitle", "Year": "TestYear", "Genre": "Test, Test, Test", "Rate": 10}

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
	key := vars["id"]
	for _, movie := range movies {
		if movie.ID == key {
			json.NewEncoder(w).Encode(movie)
		}
	}
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Extract the `id` of the movie we wish to delete
	id := vars["id"]
	// Loop through all our movies
	for index, movie := range movies {
		// if our id path parameter matches one of our
		// articles
		if movie.ID == id {
			// updates our Articles array to remove the
			// article
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
	APIVER := "/api/v1"

	// - /api/v1/post - HTTP GET request - All Posts
	// - /api/v1/post/:id - HTTP GET request - Single post
	// - /api/v1/post/:id - HTTP POST request - Publish a post
	// - /api/v1/post/:id - HTTP PATCH request - update an existing post
	// - /api/v1/post/:id - HTTP DELETE request - deletes a post

	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/amin", aminPage)
	myRouter.HandleFunc("/inc", incrementCounterPage)
	myRouter.HandleFunc(fmt.Sprintf("%s/movie", APIVER), getMovies)
	// it only calls the addMovie when http method is a POST
	// NOTE: Ordering is important here! This has to be defined before
	// the other `/movie` endpoint.
	myRouter.HandleFunc(fmt.Sprintf("%s/movie", APIVER), addMovie).Methods("POST")
	myRouter.HandleFunc(fmt.Sprintf("%s/movie/{id}", APIVER), deleteMovie).Methods("DELETE")
	myRouter.HandleFunc(fmt.Sprintf("%s/movie/{id}", APIVER), updateMovie).Methods("PUT")
	myRouter.HandleFunc(fmt.Sprintf("%s/movie/{id}", APIVER), getMovieByID)
	// http.Handle("/", http.FileServer(http.Dir("./static")))
	myRouter.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	// http.HandleFunc("/", indexPage)

	log.Fatal(http.ListenAndServe(":12345", myRouter))

	// for https running
	// log.Fatal(http.ListenAndServeTLS(":443", "server.crt", "server.key", nil))

}

func main() {
	fmt.Println("Rest API v2.0 - with Mux Routers")
	handleRequests()
}
