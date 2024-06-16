.PHONY: all
all: format lint

format:
	golangci-lint run --fix

lint:
	golangci-lint run

fast-lint:
	golangci-lint run --fast	
