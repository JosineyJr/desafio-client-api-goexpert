package repositories

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/JosineyJr/desafio-client-api-goexpert/server/models"
)

type exchangeRepository struct {
	connection *sql.DB
}

func (repository *exchangeRepository) Insert(exchange models.Exchange) error {
	stmt, err := repository.connection.Prepare(`INSERT INTO exchanges (code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, create_date) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err = stmt.ExecContext(ctx, exchange.Code, exchange.Codein, exchange.Name, exchange.High, exchange.Low, exchange.VarBid, exchange.PctChange, exchange.Bid, exchange.Ask, exchange.Timestamp, exchange.Create_date)
	if err != nil {
		return err
	}

	return nil
}

func NewExchangeRepository(connection *sql.DB) (models.IExchangeRepository, error) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS exchanges (
		code TEXT,
		codein TEXT,
		name TEXT,
		high TEXT,
		low TEXT,
		varBid TEXT,
		pctChange TEXT,
		bid TEXT,
		ask TEXT,
		timestamp TEXT,
		create_date TEXT
	);`

	_, err := connection.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	log.Println("Tabela criada ou já existente")

	exchangeRepository := exchangeRepository{connection: connection}

	return &exchangeRepository, nil
}
