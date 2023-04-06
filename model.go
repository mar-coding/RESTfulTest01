package main

import (
	"database/sql"
	"errors"
)

type Movie struct {
	Id        int     `json:"id"`
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

func (p *Movie) getMovie(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (p *Movie) updateMovie(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (p *Movie) deleteMovie(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (p *Movie) createMovie(db *sql.DB) error {
	return errors.New("Not implemented")
}

func getMovies(db *sql.DB, start, count int) ([]Movie, error) {
	return nil, errors.New("Not implemented")
}
