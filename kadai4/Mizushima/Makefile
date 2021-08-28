BINARY_NAME := bin/omikuji
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=$(GOCMD) fmt

all: build test

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

test: build
	$(GOTEST) fortune_server_test.go fortune_server.go main_test.go main.go -v -cover

clean:
	$(GOCLEAN)

format:
	$(GOFMT) ./...

.PHONY: test clean format