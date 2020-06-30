package apify

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/trickstersio/rucksack-go/apierr"
)

// Response contains information about status, body and headers of the HTTP response
type Response struct {
	status  int
	body    []byte
	headers map[string]string
}

// Builder is responsible for modifying HTTP response
type Builder interface {
	Build(r *Response) error
}

// BuilderFunc is a function which implements Builder interface
type BuilderFunc func(r *Response) error

// Build is defined to follow the Builder interfaces and calls the function itself
func (f BuilderFunc) Build(r *Response) error {
	return f(r)
}

// Status changes status of the response to the specified one
func Status(status int) BuilderFunc {
	return func(r *Response) error {
		r.status = status
		return nil
	}
}

// Error saves specified error to the response
func Error(err *apierr.Err) BuilderFunc {
	return JSON(map[string]interface{}{
		"errors": err,
	})
}

// JSON saves specified body in form of JSON object to the response
func JSON(body interface{}) BuilderFunc {
	return func(r *Response) error {
		data, err := json.Marshal(body)

		if err != nil {
			return fmt.Errorf("failed to write json response: %w", err)
		}

		r.headers["Content-Type"] = "application/json"
		r.body = data

		return nil
	}
}

// Respond sends respond built using specified builders to the client
func Respond(w http.ResponseWriter, builders ...Builder) {
	response := &Response{
		status:  http.StatusOK,
		headers: map[string]string{},
	}

	for _, builder := range builders {
		if err := builder.Build(response); err != nil {
			log.Println("Failed to build response", err)

			response.body = nil
			response.headers = map[string]string{}
			response.status = http.StatusInternalServerError
		}
	}

	for k, v := range response.headers {
		w.Header().Set(k, v)
	}

	w.WriteHeader(response.status)

	if response.body != nil {
		if _, err := w.Write(response.body); err != nil {
			log.Println("Failed to write response body", err)
		}
	}
}
