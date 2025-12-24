package errors

import "fmt"

type SessionNotFoundError struct {
	Message string
}

func (e SessionNotFoundError) Error() string {
	return fmt.Sprintf("SessionNotFound error: %s", e.Message)
}
