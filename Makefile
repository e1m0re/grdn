lint:
	golangci-lint run ./...

format-code:
	goimports --local "github.com/e1m0re/grdn" -w .

test:
	go test -race -covermode=atomic -coverprofile=coverage.out ./...

run-agent:
	go run cmd/agent/main.go

build-agent:
	go build -o bin/agent cmd/agent/*.go

build-server:
	go build -o bin/server cmd/server/*go

build:
	go build -o bin/agent cmd/agent/*.go
	go build -o bin/server cmd/server/*.go

run-server:
	go run cmd/server/main.go

run-server-db:
	DATABASE_DSN="postgresql://grdnuser:grdnpassword@192.168.33.26:5432/grdn" go run cmd/server/main.go

doc:
	godoc -http=:8081 -goroot="/Users/elmore/go/src/grdn/"
