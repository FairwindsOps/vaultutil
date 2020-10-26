# Go parameters
GOCMD=GO111MODULE=on go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

all: lint test
lint:
	golangci-lint run
reportcard:
	goreportcard-cli -t 100 -v
test:
	GO111MODULE=on $(GOCMD) test -v --bench --benchmem -coverprofile coverage.txt -covermode=atomic ./...
	GO111MODULE=on $(GOCMD) vet ./... 2> govet-report.out
	GO111MODULE=on $(GOCMD) tool cover -html=coverage.txt -o cover-report.html
	printf "\nCoverage report available at cover-report.html\n\n"
tidy:
	$(GOCMD) mod tidy
clean:
	$(GOCLEAN)
	$(GOCMD) fmt ./...
