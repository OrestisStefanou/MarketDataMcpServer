package errors

import (
	"fmt"
)

// HTTPError represents an error that occurs when there is an HTTP-related issue.
type HTTPError struct {
	StatusCode int
	Message    string
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("HTTP error: %s (status code: %d)", e.Message, e.StatusCode)
}

// JSONMarshalError represents an error during JSON marshalling.
type JSONMarshalError struct {
	Message string
	Err     error
}

func (e JSONMarshalError) Error() string {
	return fmt.Sprintf("JSON marshal error: %s: %v", e.Message, e.Err)
}

// StreamError represents an error that occurs when reading the stream from the response.
type StreamError struct {
	Message string
	Err     error
}

func (e StreamError) Error() string {
	return fmt.Sprintf("Stream error: %s: %v", e.Message, e.Err)
}
