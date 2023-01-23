# Compilation variables
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
LAMBDA=serverless_gin_lambda
PROJECT=serverless_gin_api

all: clean build package test

build:
	$(GOBUILD) -o $(LAMBDA)

package:
	zip function.zip $(LAMBDA)/*.go

test:
	$(GOTEST) -v ./...

clean:
	rm -f $(LAMBDA).zip
