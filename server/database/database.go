package database

import (
	"database/sql"

	"github.com/JosineyJr/desafio-client-api-goexpert/server/database/repositories"
	"github.com/JosineyJr/desafio-client-api-goexpert/server/models"
	_ "github.com/mattn/go-sqlite3"
)

type entityRepositories struct {
	Connection *sql.DB
	Exchange   models.IExchangeRepository
}

func Setup() (*entityRepositories, error) {
	connection, err := newSQLiteConnection()
	if err != nil {
		return nil, err
	}

	exchangeRepository, err := repositories.NewExchangeRepository(connection)
	if err != nil {
		return nil, err
	}

	return &entityRepositories{Exchange: exchangeRepository, Connection: connection}, nil
}

func newSQLiteConnection() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./database/exchange.db")
	if err != nil {
		return nil, err
	}

	return db, nil
}
