package services

import (
	"errors"
	"market_data_mcp_server/pkg/domain"
	market_data_mcp_serverErr "market_data_mcp_server/pkg/errors"
	"time"
)

type UserContextRepository interface {
	GetUserContext(userID string) (domain.UserContext, error)
	InsertUserContext(userContext domain.UserContext) error
	UpdateUserContext(userContext domain.UserContext) error
}

type UserContextDataService interface {
	GetUserContext(userID string) (domain.UserContext, error)
}

type UserContextService struct {
	userContextRepository UserContextRepository
}

func NewUserContextService(userContextRepository UserContextRepository) (*UserContextService, error) {
	return &UserContextService{userContextRepository: userContextRepository}, nil
}

func (s *UserContextService) GetUserContext(userID string) (domain.UserContext, error) {
	userContext, err := s.userContextRepository.GetUserContext(userID)
	if err != nil {
		return domain.UserContext{}, err
	}

	return userContext, nil
}

func (s *UserContextService) CreateUserContext(userContext domain.UserContext) error {
	// Check if user context for given user id already exists
	dbUserContext, err := s.userContextRepository.GetUserContext(userContext.UserID)
	if err != nil {
		notFoundError := market_data_mcp_serverErr.UserContextNotFoundError{UserID: dbUserContext.UserID}
		if errors.As(err, &notFoundError) {
			// do nothing in this case
		} else {
			return err
		}
	}

	if dbUserContext.UserID != "" {
		return market_data_mcp_serverErr.UserContextAlreadyExistsError{UserID: dbUserContext.UserID}
	}

	userContext.CreatedAt = time.Now().Format(time.RFC3339)

	return s.userContextRepository.InsertUserContext(userContext)
}

func (s *UserContextService) UpdateUserContext(userContext domain.UserContext) error {
	dbUserContext, err := s.userContextRepository.GetUserContext(userContext.UserID)
	if err != nil { // user context not found error is covered here
		return err
	}

	userContext.CreatedAt = dbUserContext.CreatedAt
	userContext.UpdatedAt = time.Now().Format(time.RFC3339)
	return s.userContextRepository.UpdateUserContext(userContext)
}
