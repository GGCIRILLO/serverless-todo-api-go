package todo

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

var ErrTodoNotFound = errors.New("todo not found")

type STATUS string

const (
	TODO        STATUS = "TODO"
	IN_PROGRESS STATUS = "IN_PROGRESS"
	DONE        STATUS = "DONE"
)

type ItemToDo struct {
	Pk          string  `json:"pk" dynamodbav:"pk"` // Partition Key per DynamoDB, usata per raggruppare i todo per utente
	Sk          string  `json:"sk" dynamodbav:"sk"` // Sort Key per DynamoDB, usata per identificare univocamente ogni todo all'interno del gruppo dell'utente
	Title       string  `json:"title" dynamodbav:"title"` // mapping per il Json e per DynamoDB: come devono chiamarsi
	Description *string `json:"description,omitempty" dynamodbav:"description,omitempty"`
	Status      STATUS  `json:"status" dynamodbav:"status"`
	CreatedAt   string  `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt   string  `json:"updatedAt" dynamodbav:"updatedAt"`
}

func (item *ItemToDo) PrepareForCreate() {
	now := time.Now().Format(time.RFC3339)
	randomID := rand.Intn(1000000)
	item.Pk = "USER#demo"
	item.Sk = fmt.Sprintf("TODO#%d", randomID)
	item.CreatedAt = now
	item.UpdatedAt = now
}

func (item *ItemToDo) PrepareForUpdate() {
	item.UpdatedAt = time.Now().Format(time.RFC3339)
}

func ValidateStatus(status string) bool {
	switch status {
	case string(TODO), string(IN_PROGRESS), string(DONE):
		return true
	default:
		return false
	}
}

func ValidateTitle(title string) bool {
	return len(title) > 0
}

func ValidateDescription(description *string) bool {
	if description == nil {
		return true
	}
	return len(*description) > 0
}

func ValidateItem(item ItemToDo) bool {
	return ValidateTitle(item.Title) && ValidateStatus(string(item.Status)) && ValidateDescription(item.Description)
}
