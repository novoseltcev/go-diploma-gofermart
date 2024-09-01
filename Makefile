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

DATABASE_URI=postgresql://postgres:postgres@0.0.0.0:5432/praktikum?sslmode=disable
.PHONY: gophermart
gophermart: $(CMD)
	DATABASE_URI=$(DATABASE_URI) JWT_LIFETIME=1 JWT_SECRET=secret123 $(CMD) -a :8080

up:
	docker-compose up -d --build

down:
	docker-compose down

.PHONY: psql
psql:
	psql $(DATABASE_URI)

PLATFORM=darwin_arm64
.PHONY: accural
accural: ./cmd/accrual/accrual_$(PLATFORM)
	./cmd/accrual/accrual_$(PLATFORM)

migrate:
	migrate -source file://migrations -database $(DATABASE_URI) up

MOCKS_DESTINATION=mocks
.PHONY: mocks
# put the files with interfaces you'd like to mock in prerequisites
# wildcards are allowed
mocks: gophermart/domains/*/service.go
	@echo "Generating mocks..."
	@rm -rf $(MOCKS_DESTINATION)
	@for file in $^; do mockgen -source=$$file -destination=$(MOCKS_DESTINATION)/$$file; done
