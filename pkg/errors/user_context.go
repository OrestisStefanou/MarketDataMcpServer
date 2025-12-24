package errors

import "fmt"

type UserContextNotFoundError struct {
	UserID string
}

func (e UserContextNotFoundError) Error() string {
	return fmt.Sprintf("user context not found for user_id: %s", e.UserID)
}

type UserContextAlreadyExistsError struct {
	UserID string
}

func (e UserContextAlreadyExistsError) Error() string {
	return fmt.Sprintf("user context for user_id: %s already exists", e.UserID)
}
