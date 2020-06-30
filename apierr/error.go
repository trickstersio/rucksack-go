package apierr

import (
	"strconv"
)

// Err contains information about object related errors together with properties related errors
type Err struct {
	Messages []string        `json:"messages,omitempty"`
	Fields   map[string]*Err `json:"attributes,omitempty"`
}

// New creates new error with object related messages
func New(messages ...string) *Err {
	return &Err{
		Messages: messages,
		Fields:   map[string]*Err{},
	}
}

// Error creates error with list of errros as object related messages
func Error(errors ...error) *Err {
	messages := make([]string, len(errors))

	for i := range errors {
		messages[i] = errors[i].Error()
	}

	return New(messages...)
}

// Field creates new error with field specific message
func Field(name string, fieldErr *Err) *Err {
	err := New()
	err.Fields[name] = fieldErr
	return err
}

// Array create new error with array item specific message
func Array(index int, itemError *Err) *Err {
	err := New()
	err.Fields[strconv.Itoa(index)] = itemError
	return err
}

// Add is adding one error data into another error
func (err *Err) Add(otherErr *Err) *Err {
	err.Messages = append(err.Messages, otherErr.Messages...)

	for fieldName, fieldErr := range otherErr.Fields {
		if existingFieldErr, ok := err.Fields[fieldName]; ok {
			existingFieldErr.Add(fieldErr)
		} else {
			err.Fields[fieldName] = fieldErr
		}
	}

	return err
}
