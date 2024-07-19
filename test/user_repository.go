package main

import (
	"context"
	"database/sql"
)

type UserRepository interface {
	Create(ctx context.Context, user User) error
}

type userRepository struct {
	db Querier
}

func NewUserRepository(db Querier) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user User) error {
	querier := r.db

	tx, ok := ctx.Value("tx").(*sql.Tx)
	if ok {
		querier = tx
	}

	_, err := querier.Exec("INSERT INTO users (name, email) VALUES ($1, $2)", user.Name, user.Email)
	if err != nil {
		return err
	}

	return nil
}
