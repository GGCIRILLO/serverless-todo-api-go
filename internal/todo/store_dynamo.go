package todo

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoStore struct {
	client    *dynamodb.Client
	tableName string
}

// usata in main.go per inizializzare il DynamoStore con il client DynamoDB e il nome della tabella
func NewDynamoStore(client *dynamodb.Client, tableName string) *DynamoStore {
	return &DynamoStore{
		client:    client,
		tableName: tableName,
	}
}

func (ds *DynamoStore) CreateTodo(ctx context.Context, todo ItemToDo) error {
	// Marshall del todo item in una mappa di attributi per DynamoDB
	// serve per convertire la struct ItemToDo in un formato che DynamoDB può comprendere
	av, err := attributevalue.MarshalMap(todo)
	if err != nil {
		return err
	}

	// Inserimento del todo item nel DynamoDB usando PutItem
	// Nella chiamata uso &dynamodb.PutItemInput: puntatore a una struct che contiene i parametri per l'operazione PutItem,
	// tra cui il nome della tabella e la mappa di attributi del todo item
	_, err = ds.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &ds.tableName,
		Item:      av,
	})
	return err
}

func (ds *DynamoStore) GetTodo(ctx context.Context, id string) (*ItemToDo, error) {
	// Definizione della chiave primaria (Partition Key ed eventualmente Sort Key)
	key, err := attributevalue.MarshalMap(map[string]string{
		"pk": "USER#demo",
		"sk": "TODO#" + id,
	})
	if err != nil {
		return nil, err
	}

	// Recupero dell'elemento usando GetItem
	result, err := ds.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &ds.tableName,
		Key:       key,
	})
	if err != nil {
		return nil, err
	}

	// Se l'elemento non viene trovato, torniamo nil senza errore
	if result.Item == nil {
		return nil, nil
	}

	// Conversione dei dati da DynamoDB alla nostra struct Go
	var todo ItemToDo
	err = attributevalue.UnmarshalMap(result.Item, &todo)
	if err != nil {
		return nil, err
	}

	return &todo, nil
}


func (ds *DynamoStore) ListToDosByUser(ctx context.Context, userID string) ([]ItemToDo, error) {
	out, err := ds.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              &ds.tableName,
		KeyConditionExpression: aws.String("pk = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: "USER#" + userID},
		},
	})
	if err != nil {
		return nil, err
	}

	var todos []ItemToDo
	err = attributevalue.UnmarshalListOfMaps(out.Items, &todos)
	if err != nil {
		return nil, err
	}
	return todos, nil
}

func (ds *DynamoStore) UpdateTodo(ctx context.Context, todo ItemToDo) error {
	// Marshall dell'intero oggetto aggiornato
	av, err := attributevalue.MarshalMap(todo)
	if err != nil {
		return err
	}

	// UpdateTodo usa lo stesso PutItem: in DynamoDB PutItem sovrascrive interamente l'elemento
	// se la chiave primaria (pk e sk) coincide
	_, err = ds.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &ds.tableName,
		Item:      av,
	})
	return err
}

func (ds *DynamoStore) DeleteTodo(ctx context.Context, id string) error {
	// Definizione della chiave dell'elemento da eliminare
	key, err := attributevalue.MarshalMap(map[string]string{
		"pk": "USER#demo",
		"sk": "TODO#" + id,
	})
	if err != nil {
		return err
	}

	// Esecuzione dell'operazione di eliminazione
	_, err = ds.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: &ds.tableName,
		Key:       key,
	})
	return err
}
