package main

import (
	"context"
	"uow"
)

type SomeService interface {
	DoSomething(ctx context.Context, user User, movie Movie) error
}

type someService struct {
	userRepo  UserRepository
	movieRepo MovieRepository
	uowBase   uow.Base
}

func NewSomeService(
	userRepo UserRepository,
	movieRepo MovieRepository,
	uowBase uow.Base,
) SomeService {
	return &someService{
		userRepo:  userRepo,
		movieRepo: movieRepo,
		uowBase:   uowBase,
	}
}

func (s *someService) DoSomething(ctx context.Context, user User, movie Movie) (err error) {
	uowInstance := uow.NewInstance(s.uowBase)

	ctx, err = uowInstance.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = uowInstance.Rollback()
		}
	}()

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return err
	}

	err = s.movieRepo.Create(ctx, movie)
	if err != nil {
		return err
	}

	err = uowInstance.Commit()
	if err != nil {
		return err
	}

	return nil
}
