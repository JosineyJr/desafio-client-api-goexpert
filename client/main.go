package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type exchangeGatewayResponse struct {
	Bid string `json:"bid"`
}

func createFile(filePath string) (*os.File, error) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func getExchange() (*exchangeGatewayResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:3000/cotacao", nil)
	if err != nil {
		return nil, err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var exchange exchangeGatewayResponse
	json.Unmarshal(body, &exchange)

	return &exchange, nil
}

func main() {
	log.Println("Buscando cotação...")
	exchange, err := getExchange()
	if err != nil {
		panic(err)
	}
	log.Printf("Cotação: %v\n", exchange.Bid)

	file, err := createFile("./reports/cotacao.txt")
	if err != nil {
		panic(err)
	}

	fmt.Println(time.Local)

	_, err = file.WriteString(fmt.Sprintf("Data: %v | Dólar: %v\n", time.Now().Format("2006-01-02 15:04:05"), exchange.Bid))
	if err != nil {
		panic(err)
	}
	log.Println("Cotação atual salva no arquivo")
}
