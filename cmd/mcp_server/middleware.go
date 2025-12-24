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
			m.logger.Printf("Tool call failed: session=%s tool=%s duration=%v error=%v",
				"0", req.Params.Name, duration, err)
		} else {
			m.logger.Printf("Tool call completed: session=%s tool=%s duration=%v",
				"0", req.Params.Name, duration)
		}

		return result, err
	}
}
