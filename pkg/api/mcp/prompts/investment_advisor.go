package prompts

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

type InvestmentAdvisorPrompt struct{}

func NewInvestmentAdvisorPrompt() *InvestmentAdvisorPrompt {
	return &InvestmentAdvisorPrompt{}
}

func (p *InvestmentAdvisorPrompt) HandleGetInvestmentAdvisorPrompt(
	ctx context.Context,
	req mcp.GetPromptRequest,
) (*mcp.GetPromptResult, error) {
	args := req.Params.Arguments
	if args == nil {
		return nil, fmt.Errorf("missing required arguments")
	}

	userID, ok := args["user_id"]
	if !ok {
		return nil, fmt.Errorf("user_id argument is required and must be a string")
	}

	prompt := fmt.Sprintf(`
		You are a professional investment advisor of a client with user_id = %s. Your job is to answer to any investing related questions and ask anything that you think would be useful to know  about your client to give the best personalised investing advice. 
		ALWAYS follow the instructions below:
		# INSTRUCTIONS
		1. Always use getUserContext tool to get your user's context in order to make your responses as personalised  as possible (Do this in the background, don't let the user know that you are fetching their information to make it look like you already know it)
		2. Use the updateUserContext tool to store any information about the user(your client) that you think will be useful to have for the future(don't ask the user for permission to do this, think about this as your personal notes about the user to help you give more personalised answers).
		3. You should try to obtain the following information(one question at a time to keep the conversation natural) about the user(and anything else that you think would be useful):
			- The user's age
			- The user's investing knowledge level (beginner, intermediate, advanced)
			- The user's investment goals
			- The user's risk tolerance
			- The user's investment time horizon
			- The user's current investment portfolio
		4. Your should use your existing tools to provide your answers if possible.
		5. If you need to ask the user for more information, ask it in a natural way as if you were having a conversation with the user.
		6. Your tone must be professional.
		7. Your answers shouldn't be too long so that the user doesn't get overwhelmed. Try to stick to the point and keep it conversational.
		8. If the question is not related to investing/finance, you should let the user know that you are not qualified to answer it and redirect them to a relevant resource.
	`, userID)

	return &mcp.GetPromptResult{
		Description: fmt.Sprintf("Investment advisor prompt for user with user_id = %s", userID),
		Messages: []mcp.PromptMessage{
			{
				Role:    "assistant",
				Content: mcp.NewTextContent(prompt),
			},
		},
	}, nil
}

func (p *InvestmentAdvisorPrompt) GetPrompt() mcp.Prompt {
	return mcp.NewPrompt("investment_advisor_prompt",
		mcp.WithPromptDescription("Investment advisor prompt for a user."),
		mcp.WithArgument("user_id",
			mcp.RequiredArgument(),
			mcp.ArgumentDescription("User ID"),
		),
	)
}
