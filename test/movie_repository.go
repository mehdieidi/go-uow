package main

import (
	"context"
	"database/sql"
)

type MovieRepository interface {
	Create(ctx context.Context, Movie Movie) error
}

type movieRepository struct {
	db Querier
}

func NewMovieRepository(db Querier) MovieRepository {
	return &movieRepository{db: db}
}

func (r *movieRepository) Create(ctx context.Context, m Movie) error {
	querier := r.db

	tx, ok := ctx.Value("tx").(*sql.Tx)
	if ok {
		querier = tx
	}

	_, err := querier.Exec("INSERT INTO movies (name) VALUES ($1)", m.Name)
	if err != nil {
		return err
	}

	return nil
}
