package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/JosineyJr/desafio-client-api-goexpert/server/database"
	"github.com/JosineyJr/desafio-client-api-goexpert/server/models"
)

type exchangeApiResponse struct {
	USDBRL models.Exchange `json:"USDBRL"`
}

type apiResponse struct {
	Bid string `json:"bid"`
}

type returnMessage struct {
	Message string `json:"message"`
}

type exchangeController struct {
	exchangeRepository models.IExchangeRepository
}

func main() {
	entityRepositories, err := database.Setup()
	if err != nil {
		panic(err)
	}
	defer entityRepositories.Connection.Close()

	startServer(&entityRepositories.Exchange)
}

func startServer(db *models.IExchangeRepository) {
	mux := http.NewServeMux()

	exchangeHandler := exchangeController{exchangeRepository: *db}

	mux.Handle("/cotacao", exchangeHandler)

	log.Println("Server started on port 3000")
	http.ListenAndServe(":3000", mux)
}

func (controller *exchangeController) insertExchangeInDb(exchange *models.Exchange) error {
	err := controller.exchangeRepository.Insert(*exchange)
	if err != nil {
		return err
	}

	log.Println("Cotação inserida no banco de dados")

	return nil
}

func (controller exchangeController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(http.StatusInternalServerError)
		jsonEncoder.Encode(returnMessage{Message: "Internal Server Error"})

		log.Printf("Erro ao capturar o response body: %v", err)
		return
	}

	var exchangeResponse exchangeApiResponse
	err = json.Unmarshal(body, &exchangeResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		jsonEncoder.Encode(returnMessage{Message: "Internal Server Error"})

		log.Printf("Erro ao desserializar json: %v", err)
		return
	}

	err = controller.insertExchangeInDb(&exchangeResponse.USDBRL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		jsonEncoder.Encode(returnMessage{Message: "Internal Server Error"})

		log.Printf("Erro ao inserir cotação no banco de dados: %v", err)
		return
	}

	log.Println("Requisição concluída")

	w.WriteHeader(http.StatusOK)
	jsonEncoder.Encode(apiResponse{Bid: exchangeResponse.USDBRL.Bid})
}
