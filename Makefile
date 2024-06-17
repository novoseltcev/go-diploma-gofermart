.PHONY: all
all: generate format lint build up gophermart

format:
	golangci-lint run --fix

lint:
	golangci-lint run

fast-lint:
	golangci-lint run --fast	

generate:
	go generate ./...

DIR=./cmd/gophermart
CMD=$(DIR)/gophermart

build: generate $(DIR)/main.go
	go build -buildvcs=false -o $(CMD) $(DIR) 

.PHONY: gophermart
gophermart: $(CMD)
	DATABASE_DSN="postgres://postgres:postgres@0.0.0.0:5432/praktikum?sslmode=disable" $(CMD) -a :8080

up:
	docker-compose up -d --build

down:
	docker-compose down

psql:
	PGPASSWORD=postgres psql -U postgres -h 0.0.0.0 -p 5432 -d praktikum

PLATFORM=darwin_arm64
.PHONY: accural
accural: ./cmd/accrual/accrual_$(PLATFORM)
	./cmd/accrual/accrual_$(PLATFORM)
