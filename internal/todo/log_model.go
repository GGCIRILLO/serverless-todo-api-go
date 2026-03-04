package todo

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
)

type logEntry struct {
	Level     string `json:"level"`
	Message   string `json:"message"`
	Method    string `json:"method,omitempty"`
	Path      string `json:"path,omitempty"`
	RequestID string `json:"requestId,omitempty"`
	Error     string `json:"error,omitempty"`
}

func LogInfo(req events.APIGatewayV2HTTPRequest, msg string) {
	e := logEntry{
		Level:     "INFO",
		Message:   msg,
		Method:    req.RequestContext.HTTP.Method,
		Path:      req.RawPath,
		RequestID: req.RequestContext.RequestID,
	}
	b, _ := json.Marshal(e)
	log.Println(string(b))
}

func LogError(req events.APIGatewayV2HTTPRequest, msg string, err error) {
	e := logEntry{
		Level:     "ERROR",
		Message:   msg,
		Method:    req.RequestContext.HTTP.Method,
		Path:      req.RawPath,
		RequestID: req.RequestContext.RequestID,
	}
	if err != nil {
		e.Error = err.Error()
	}
	b, _ := json.Marshal(e)
	log.Println(string(b))
}
