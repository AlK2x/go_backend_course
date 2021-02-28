all: build test check

modules:
	go mod tidy

build: modules
	go build -v -o bin/orderservice cmd/orderservice/*.go

test:
	go test ./...

check:
	golangci-lint run