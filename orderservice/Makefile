all: build

modules:
	go mod tidy

build:
	go build -v -o bin/orderservice cmd/orderservice/*.go
