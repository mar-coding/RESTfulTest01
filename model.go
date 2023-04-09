package main

import (
	"context"
	"database/sql"
	"fmt"
)

type Movie struct {
	Id        int64   `json:"id"`
	Name      string  `json:"name"`
	Year      string  `json:"year"`
	Genre     string  `json:"genre"`
	Duration  string  `json:"duration"`
	Origin    string  `json:"origin"`
	Director  string  `json:"director"`
	Rate      float32 `json:"rate"`
	Ratecount int     `json:"rate_count"`
	Link      string  `json:"link"`
}

func (m *Movie) getMovie(db *sql.DB) error {
	q := "SELECT * FROM Movie WHERE movie_id=?"
	res := db.QueryRow(q, m.Id).Scan(&m.Id, &m.Name, &m.Year, &m.Genre, &m.Duration, &m.Origin, &m.Director, &m.Rate, &m.Ratecount, &m.Link)
	return res
}

func (m *Movie) updateMovie(db *sql.DB) error {
	q, err := db.Prepare("UPDATE Movie SET movie_name=?, movie_year=?, movie_genre=?, movie_duration=?, movie_origin=?, movie_director=?, movie_rating=?, movie_rating_count=?, movie_link=? WHERE movie_id=?")
	q.Exec(m.Name, m.Year, m.Genre, m.Duration, m.Origin, m.Director, m.Rate, m.Ratecount, m.Link, m.Id)
	return err
}

func (m *Movie) deleteMovie(db *sql.DB) error {
	q, err := db.Prepare("DELETE FROM Movie WHERE movie_id=?")
	q.Exec(m.Id)
	return err
}

func (m *Movie) createMovie(db *sql.DB) error {
	q := fmt.Sprintf("INSERT INTO Movie(movie_id, movie_name, movie_year, movie_genre, movie_duration, movie_origin, movie_director, movie_rating, movie_rating_count, movie_link) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	insertResult, err := db.ExecContext(context.Background(), q, m.Id, m.Name, m.Year, m.Genre, m.Duration, m.Origin, m.Director, m.Rate, m.Ratecount, m.Link)
	if err != nil {
		return err
	}
	m.Id, err = insertResult.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}

func getMovies(db *sql.DB, start, count int) ([]Movie, error) {
	q := fmt.Sprintf("SELECT * FROM Movie ORDER BY movie_id ASC LIMIT %d OFFSET %d", count, start)
	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	movies := []Movie{}
	for rows.Next() {
		var m Movie
		if err := rows.Scan(&m.Id, &m.Name, &m.Year, &m.Genre, &m.Duration, &m.Origin, &m.Director, &m.Rate, &m.Ratecount, &m.Link); err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}

	return movies, nil
}
