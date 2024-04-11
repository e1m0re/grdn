lint:
	golangci-lint run ./...

test:
	go test -race -covermode=atomic -coverprofile=coverage.out ./...

run-agent:
	go run cmd/agent/main.go cmd/agent/flags.go

run-server:
	go run cmd/server/main.go

run-server-db:
	DATABASE_DSN="postgresql://grdnuser:grdnpassword@192.168.33.26:5432/grdn" go run cmd/server/main.go