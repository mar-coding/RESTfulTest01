package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// movie represents data about a movie.
type movie struct {
	ID    string  `json:"id"`
	Title string  `json:"title"`
	Year  string  `json:"year"`
	Genre string  `json:"genre"`
	Rate  float32 `json:"rate"`
}

// movies slice to seed movie data.
var movies = []movie{
	{ID: "0", Title: "The Shawshank Redemption", Year: "1994", Genre: "Drama", Rate: 9.3},
	{ID: "1", Title: "The Godfather", Year: "1972", Genre: "Crime, Drama", Rate: 9.2},
	{ID: "2", Title: "The Dark Knight", Year: "2008", Genre: "Action, Crime, Drama", Rate: 9.0},
	{ID: "3", Title: "The Godfather Part II", Year: "1974", Genre: "Crime, Drama", Rate: 9.0},
	{ID: "4", Title: "Schindler's List", Year: "1993", Genre: "Biography, Drama, History", Rate: 9.0},
}

// get all movies in json format
func getMovies(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, movies)
}

func postMovies(c *gin.Context) {
	var newMovie movie
	if err := c.BindJSON(&newMovie); err != nil {
		return
	}

	movies = append(movies, newMovie)
	c.IndentedJSON(http.StatusCreated, newMovie)
}

func getMovieByID(c *gin.Context) {
	id := c.Param("id")
	for _, a := range movies {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "movie not found"})
}

func main() {
	router := gin.Default()
	router.GET("/movies", getMovies)
	router.GET("/movies/:id", getMovieByID)
	router.POST("/movies", postMovies)
	router.Run("localhost:8080")
}
