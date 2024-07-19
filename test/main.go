package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"strconv"
	"uow"
)

func main() {
	fmt.Println("[+] Testing...")

	postgresDB, err := sql.Open(
		"postgres",
		fmt.Sprintf(
			"host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
			"localhost",
			"5432",
			"postgres",
			"1234",
			"uow_test",
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	err = postgresDB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	userRepo := NewUserRepository(postgresDB)
	movieRepo := NewMovieRepository(postgresDB)

	uowBase := uow.NewBase(postgresDB)

	someService := NewSomeService(
		userRepo,
		movieRepo,
		uowBase,
	)

	for i := 0; i < 100; i++ {
		fmt.Println("Iteration", strconv.Itoa(i))

		userSample := User{
			Name:  "some user" + strconv.Itoa(i),
			Email: "someuser@example.com" + strconv.Itoa(i),
		}
		movieSample := Movie{
			Name: "some movie" + strconv.Itoa(i),
		}

		err = someService.DoSomething(context.Background(), userSample, movieSample)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("[-] Finished")
}
