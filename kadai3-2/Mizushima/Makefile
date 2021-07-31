BINARY_NAME := bin/paraDW
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
	cd download; $(GOTEST) -v -cover
	cd getheader; $(GOTEST) -v -cover
	cd listen; $(GOTEST) -v -cover
	cd request; $(GOTEST) -v -cover

clean:
	$(GOCLEAN)

format:
	$(GOFMT) ./...

.PHONY: test clean format
