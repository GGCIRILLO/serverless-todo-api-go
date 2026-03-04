package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/GGCIRILLO/serverless-todo-api-go/internal/todo"
)

// Dichiaro un puntatore globale al DynamoStore per utilizzarlo poi nel handler della lambda
var store *todo.DynamoStore

func init() {
	tableName := os.Getenv("TABLE_NAME")
	if tableName == "" {
		tableName = "todos" // Default per sviluppo
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	dbClient := dynamodb.NewFromConfig(cfg)
	store = todo.NewDynamoStore(dbClient, tableName)
}

type healthResponse struct {
	Status string `json:"status"`
}

type errorResponse struct {
	Message string `json:"message"`
}

func apiResponse(status int, body interface{}) events.APIGatewayV2HTTPResponse {
	stringBody := ""
	if body != nil {
		b, _ := json.Marshal(body)
		stringBody = string(b)
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: status,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       stringBody,
	}
}

func handleGetTodos(ctx context.Context) events.APIGatewayV2HTTPResponse {
	todos, err := store.ListToDosByUser(ctx, "demo")
	if err != nil {
		return apiResponse(http.StatusInternalServerError, errorResponse{Message: "error fetching todos"})
	}
	return apiResponse(http.StatusOK, todos)
}

func handleCreateTodo(ctx context.Context, body string) events.APIGatewayV2HTTPResponse {
	var item todo.ItemToDo
	if err := json.Unmarshal([]byte(body), &item); err != nil {
		return apiResponse(http.StatusBadRequest, errorResponse{Message: "invalid request body"})
	}

	// Popolamento automatico campi
	item.PrepareForCreate()

	if err := store.CreateTodo(ctx, item); err != nil {
		return apiResponse(http.StatusInternalServerError, errorResponse{Message: "failed to create todo"})
	}
	return apiResponse(http.StatusCreated, item)
}

func handleGetTodo(ctx context.Context, request events.APIGatewayV2HTTPRequest) events.APIGatewayV2HTTPResponse {
	id := request.PathParameters["id"]
	if id == "" {
		return apiResponse(http.StatusBadRequest, errorResponse{Message: "missing id parameter"})
	}

	item, err := store.GetTodo(ctx, id)
	if err != nil {
		return apiResponse(http.StatusInternalServerError, errorResponse{Message: "failed to get todo"})
	}
	if item == nil {
		return apiResponse(http.StatusNotFound, errorResponse{Message: "todo not found"})
	}

	return apiResponse(http.StatusOK, item)
}

func handleUpdateTodo(ctx context.Context, body string) events.APIGatewayV2HTTPResponse {
	var item todo.ItemToDo
	if err := json.Unmarshal([]byte(body), &item); err != nil {
		return apiResponse(http.StatusBadRequest, errorResponse{Message: "invalid request body"})
	}

	// 1. Recuperiamo il record esistente tramite l'ID contenuto nella SK (formato "TODO#123")
	if len(item.Sk) <= 5 {
		return apiResponse(http.StatusBadRequest, errorResponse{Message: "invalid sk format"})
	}
	id := item.Sk[5:]
	existing, err := store.GetTodo(ctx, id)
	if err != nil {
		return apiResponse(http.StatusInternalServerError, errorResponse{Message: "failed to retrieve existing record"})
	}
	if existing == nil {
		return apiResponse(http.StatusNotFound, errorResponse{Message: "todo not found"})
	}

	// 2. Aggiornamento automatico UpdatedAt
	item.PrepareForUpdate()

	// 3. Eseguiamo l'update passando sia il nuovo item che quello esistente per preservare i dati
	if err := store.UpdateTodo(ctx, item, *existing); err != nil {
		return apiResponse(http.StatusInternalServerError, errorResponse{Message: "failed to update todo"})
	}

	// Restituiamo l'oggetto aggiornato (con il createdAt originale recuperato)
	item.CreatedAt = existing.CreatedAt
	return apiResponse(http.StatusOK, item)
}

func handleDeleteTodo(ctx context.Context, request events.APIGatewayV2HTTPRequest) events.APIGatewayV2HTTPResponse {
	id := request.PathParameters["id"]
	if id == "" {
		return apiResponse(http.StatusBadRequest, errorResponse{Message: "missing id parameter"})
	}

	if err := store.DeleteTodo(ctx, id); err != nil {
		return apiResponse(http.StatusInternalServerError, errorResponse{Message: "failed to delete todo"})
	}
	return apiResponse(http.StatusNoContent, nil)
}

func handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	method := request.RequestContext.HTTP.Method
	path := request.RawPath

	switch {
	case method == "GET" && path == "/health":
		return apiResponse(http.StatusOK, healthResponse{Status: "ok"}), nil

	case method == "GET" && path == "/todos":
		return handleGetTodos(ctx), nil

	case method == "POST" && path == "/todos":
		return handleCreateTodo(ctx, request.Body), nil

	case method == "GET" && (len(path) > 7 && path[:7] == "/todos/"):
		return handleGetTodo(ctx, request), nil

	case method == "PUT" && path == "/todos":
		return handleUpdateTodo(ctx, request.Body), nil

	case method == "DELETE" && (len(path) > 7 && path[:7] == "/todos/"):
		return handleDeleteTodo(ctx, request), nil

	default:
		return apiResponse(http.StatusNotFound, errorResponse{Message: "not found"}), nil
	}
}

func main() {
	// Avvio della lambda con il handler definito
	// La funzione lambda viene eseguita quando viene invocata, e il handler gestisce le richieste in arrivo
	// In questo caso, il handler risponde a una richiesta GET al percorso /health con un messaggio di stato
	// Se la richiesta non corrisponde a questo percorso, restituisce un messaggio di "not found"
	lambda.Start(handler)
}
