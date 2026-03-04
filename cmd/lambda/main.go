package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type healthResponse struct {
	Status string `json:"status"`
}
func handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	if request.RequestContext.HTTP.Method == "GET" && request.RawPath == "/health" {
		responseBody, err := json.Marshal(healthResponse{Status: "ok"})
		if err != nil {
			return events.APIGatewayV2HTTPResponse{
				StatusCode: http.StatusInternalServerError,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       `{"message":"internal server error"}`,
			}, nil
		}
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusOK,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       string(responseBody),
		}, nil
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusNotFound,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       `{"message":"not found"}`,
	}, nil
}

func main() {
	// Avvio della lambda con il handler definito
	// La funzione lambda viene eseguita quando viene invocata, e il handler gestisce le richieste in arrivo
	// In questo caso, il handler risponde a una richiesta GET al percorso /health con un messaggio di stato
	// Se la richiesta non corrisponde a questo percorso, restituisce un messaggio di "not found"
	lambda.Start(handler)
}