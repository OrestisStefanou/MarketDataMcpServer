package errors

import "fmt"

type SuperInvestorPortfolioNotFoundError struct {
	Message string
}

func (e SuperInvestorPortfolioNotFoundError) Error() string {
	return fmt.Sprintf("SuperInvestorPortfolioNotFoundError error: %s", e.Message)
}
