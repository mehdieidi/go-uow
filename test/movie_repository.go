package main

import (
	"context"
	"uow"
)

type MovieRepository interface {
	Create(ctx context.Context, Movie Movie) error
}

type movieRepository struct {
	querier Querier
}

func NewMovieRepository(querier Querier) MovieRepository {
	return &movieRepository{querier: querier}
}

func (r *movieRepository) Create(ctx context.Context, m Movie) error {
	querier := r.querier

	tx, ok := uow.TxFromContext(ctx)
	if ok {
		querier = tx
	}

	_, err := querier.Exec("INSERT INTO movies (name) VALUES ($1)", m.Name)
	if err != nil {
		return err
	}

	return nil
}
