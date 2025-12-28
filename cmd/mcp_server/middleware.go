package main

import (
	"context"
	"log"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type LoggingMiddleware struct {
	logger *log.Logger
}

func NewLoggingMiddleware(logger *log.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{logger: logger}
}

func (m *LoggingMiddleware) ToolMiddleware(next server.ToolHandlerFunc) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		start := time.Now()
		// sessionID := server.GetSessionID(ctx)	// TODO: ADD session at some point

		m.logger.Printf("Tool call started: tool=%s with args=%v", req.Params.Name, req.Params.Arguments)

		result, err := next(ctx, req)

		duration := time.Since(start)

		if err != nil {
			m.logger.Printf("Tool call failed (protocol error): tool=%s duration=%v error=%v",
				req.Params.Name, duration, err)
		} else if result.IsError {
			var errorMsg string
			if len(result.Content) > 0 {
				if textContent, ok := mcp.AsTextContent(result.Content[0]); ok {
					errorMsg = textContent.Text
				} else {
					errorMsg = "Non-text content in error result"
				}
			} else {
				errorMsg = "Empty error result content"
			}
			m.logger.Printf("Tool call failed (tool error): tool=%s args=%v duration=%v error=%s",
				req.Params.Name, req.Params.Arguments, duration, errorMsg)
		} else {
			m.logger.Printf("Tool call completed (success): tool=%s args=%v duration=%v",
				req.Params.Name, req.Params.Arguments, duration)
		}

		return result, err
	}
}
