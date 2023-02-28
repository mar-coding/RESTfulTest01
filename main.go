package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"strconv"
	"sync"
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

// get all movies in json format
func getMovies(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(movies)
}

// func postMovies(c *gin.Context) {
// 	var newMovie movie
// 	if err := c.BindJSON(&newMovie); err != nil {
// 		return
// 	}

// 	movies = append(movies, newMovie)
// 	c.IndentedJSON(http.StatusCreated, newMovie)
// }

// func getMovieByID(c *gin.Context) {
// 	id := c.Param("id")
// 	for _, a := range movies {
// 		if a.ID == id {
// 			c.IndentedJSON(http.StatusOK, a)
// 			return
// 		}
// 	}
// 	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "movie not found"})
// }

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

	http.HandleFunc("/amin", aminPage)
	http.HandleFunc("/inc", incrementCounterPage)
	http.HandleFunc("/movie", getMovies)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	// http.HandleFunc("/", indexPage)

	log.Fatal(http.ListenAndServe(":12345", nil))

	// for https running
	// log.Fatal(http.ListenAndServeTLS(":443", "server.crt", "server.key", nil))
}

func main() {
	handleRequests()
}
