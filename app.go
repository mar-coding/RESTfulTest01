package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"github.com/gorilla/mux"
)

var PORT = "9990"

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(db_name string, dsn string) {
	a.DB = dbConnect(db_name, dsn)
	a.Router = mux.NewRouter()
	a.handleRequests()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", PORT), a.Router))
}

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

// ---------------------------------------------
// ----------- Routes & API handling -----------
// ---------------------------------------------

func (a *App) apiDeleteMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Movie ID")
		return
	}

	p := Movie{Id: int64(id)}
	if err := p.deleteMovie(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) apiUpdateMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Movie ID")
		return
	}

	var m Movie
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	m.Id = int64(id)

	if err := m.updateMovie(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, m)
}

func (a *App) apiCreateMovie(w http.ResponseWriter, r *http.Request) {
	var m Movie
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := m.createMovie(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, m)
}

func (a *App) apiGetMovies(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	// fixing problem with thresholds
	MAX_COUNT := 10
	MIN_COUNT := 1
	DEFAULT_START := 0
	if count > MAX_COUNT || count < MIN_COUNT {
		count = MAX_COUNT
	}
	if start < DEFAULT_START {
		start = DEFAULT_START
	}

	movies, err := getMovies(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, movies)
}

func (a *App) apiGetMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Movie ID")
		return
	}

	m := Movie{Id: int64(id)}
	if err := m.getMovie(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Movie not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, m)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	// process error responses
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// process normal responses
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) handleRequests() {

	//API version variable
	API_VER := "/api/v1"

	// The "{id:[0-9]+}" part of the path indicates that Gorilla Mux should treat process a URL only if the id is a number.
	// For all matching requests, Gorilla Mux then stores the actual numeric value in the "id" variable

	a.Router.HandleFunc(fmt.Sprintf("%s/movie/{id:[0-9]+}", API_VER), a.apiGetMovie).Methods("GET")
	a.Router.HandleFunc(fmt.Sprintf("%s/movies", API_VER), a.apiGetMovies).Methods("GET")
	a.Router.HandleFunc(fmt.Sprintf("%s/movie", API_VER), a.apiCreateMovie).Methods("POST")
	a.Router.HandleFunc(fmt.Sprintf("%s/movie/{id:[0-9]+}", API_VER), a.apiUpdateMovie).Methods("PUT")
	a.Router.HandleFunc(fmt.Sprintf("%s/movie/{id:[0-9]+}", API_VER), a.apiDeleteMovie).Methods("DELETE")
}
