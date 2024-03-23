package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type exchange struct {
	code        string
	codein      string
	name        string
	high        string
	low         string
	varBid      string
	pctChange   string
	Bid         string `json:"bid"`
	ask         string
	timestamp   string
	create_date string
}

type apiResponse struct {
	USDBRL exchange `json:"USDBRL"`
}

type returnMessage struct {
	Message string `json:"message"`
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/cotacao", exchangeHandler)

	log.Println("Server started on port 3000")
	http.ListenAndServe(":3000", mux)
}

func exchangeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Requisição iniciada")
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	w.Header().Set("Content-Type", "application/json")

	jsonEncoder := json.NewEncoder(w)

	request, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		log.Printf("Erro ao criar request: %v", err)
		w.WriteHeader(http.StatusInternalServerError)

		jsonEncoder.Encode(returnMessage{Message: "Internal Server Error"})
		return
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Printf("Erro ao realizar requisição: %v", err)
		switch err.Error() {
		case `Get "https://economia.awesomeapi.com.br/json/last/USD-BRL": context deadline exceeded`:
			w.WriteHeader(http.StatusRequestTimeout)

			jsonEncoder.Encode(returnMessage{Message: "Exchange API timed out"})

		default:
			w.WriteHeader(http.StatusInternalServerError)

			jsonEncoder.Encode(returnMessage{Message: "Internal Server Error"})
		}

		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Erro ao capturar o response body: %v", err)
		return
	}

	var exchange apiResponse
	err = json.Unmarshal(body, &exchange)
	if err != nil {
		log.Printf("Erro ao transformar o response body para json: %v", err)
		return
	}

	log.Println("Requisição concluída")

	w.WriteHeader(http.StatusOK)
	jsonEncoder.Encode(exchange.USDBRL)
}
