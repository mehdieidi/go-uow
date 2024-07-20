# go-uow

This Go package implements the Unit of Work pattern, which helps to manage
transactions and business logic without compromising the layered
architecture. It's designed to simplify complex transaction management
in your Go applications.

## Installation

To install the package, use the following command:

```bash
$ go get github.com/mehdieidi/go-uow
```

## Usage

### Querier interface

Use querier interface in your repository layer that both the sql.DB and
sql.Tx implements.

```Go
type Querier interface {
    Query(query string, args ...any) (*sql.Rows, error)
    Exec(query string, args ...any) (sql.Result, error)
}
```

#### Repository layer

Example: user repository layer.

```Go
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
	
    // If tx exists, a transaction is started and queries must run
    // in the tx. If tx doesn't exist in the ctx, must use the general
    // repository querier.
    tx, ok := uow.TxFromContext(ctx)
    if ok {
        querier = tx
    }

    _, err := querier.Exec("INSERT INTO users (name) VALUES ($1)", user.Name)
    if err != nil {
        return err
    }

    return nil
}
```

#### Service layer

Suppose we have two repositories implemented in the format shown above.
Movie and User.

For example the service layer has a method that needs to run two queries
atomically in a transaction.

```Go
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
    tx := uow.NewTransaction(s.uowBase)

    ctx, err = tx.Begin(ctx)
    if err != nil {
        return err
    }

    defer func() {
    if err != nil {
        _ = tx.Rollback()
    }
    _ = tx.Commit()
    }()

    err = s.userRepo.Create(ctx, user)
    if err != nil {
        return err
    }

    err = s.movieRepo.Create(ctx, movie)
    if err != nil {
        return err
    }

    return nil
}
```

#### Main

You need to create a base uow using a sql.DB instance.

```Go
func main() {
    db, err := sql.Open(
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
	
    ...

    uowBase := uow.NewBase(db)
}
```

## License

MIT