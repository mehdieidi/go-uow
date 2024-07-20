package main

import (
	"context"
	"uow"
)

type UserRepository interface {
	Create(ctx context.Context, user User) error
}

type userRepository struct {
	querier Querier
}

func NewUserRepository(querier Querier) UserRepository {
	return &userRepository{querier: querier}
}

func (r *userRepository) Create(ctx context.Context, user User) error {
	querier := r.querier

	tx, ok := uow.TxFromContext(ctx)
	if ok {
		querier = tx
	}

	_, err := querier.Exec("INSERT INTO users (name, email) VALUES ($1, $2)", user.Name, user.Email)
	if err != nil {
		return err
	}

	return nil
}
