test:
	go test -race -covermode=atomic -coverprofile=coverage.out ./...

run-agent:
	go run cmd/agent/main.go cmd/agent/flags.go

run-server:
	go run cmd/server/main.go cmd/server/flags.go