package services

import (
	"encoding/json"
	"fmt"
	"market_data_mcp_server/pkg/domain"
	"os"
	"sync"
)

type InvestingIdeasService interface {
	GetInvestingIdeas() ([]domain.InvestingIdea, error)
	GetInvestingIdeaStocks(ideaID string) ([]string, error)
}

type InvestingIdeasLocalDataService struct {
	investingIdeas map[string]*investingIdeaData
	rwMutex        sync.RWMutex
}

type investingIdeaData struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Companies []string `json:"companies"`
}

func NewInvestingIdeasLocalDataService(dataPath string) (*InvestingIdeasLocalDataService, error) {
	data, err := os.ReadFile(dataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read investing ideas file: %w", err)
	}

	var ideas []investingIdeaData
	if err := json.Unmarshal(data, &ideas); err != nil {
		return nil, fmt.Errorf("failed to unmarshal investing ideas: %w", err)
	}

	investingIdeasMap := make(map[string]*investingIdeaData)
	for i := range ideas {
		investingIdeasMap[ideas[i].ID] = &ideas[i]
	}

	return &InvestingIdeasLocalDataService{
		investingIdeas: investingIdeasMap,
	}, nil
}

func (s *InvestingIdeasLocalDataService) GetInvestingIdeas() ([]domain.InvestingIdea, error) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	ideas := make([]domain.InvestingIdea, 0, len(s.investingIdeas))
	for _, data := range s.investingIdeas {
		ideas = append(ideas, domain.InvestingIdea{
			ID:    data.ID,
			Title: data.Title,
		})
	}

	return ideas, nil
}

func (s *InvestingIdeasLocalDataService) GetInvestingIdeaStocks(ideaID string) ([]string, error) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	data, ok := s.investingIdeas[ideaID]
	if !ok {
		return nil, fmt.Errorf("investing idea not found: %s", ideaID)
	}

	return data.Companies, nil
}
