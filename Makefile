run_mcp_server:
	go run ./cmd/mcp_server

build_mcp_server:
	go build cmd/mcp_server/main.go

install:
	go mod tidy
	go mod download
