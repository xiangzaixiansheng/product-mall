.PHONY: build run test lint clean docker docker-compose fmt vet

APP_NAME := product-mall
BUILD_DIR := ./bin
MAIN_PATH := ./cmd/main.go

build:
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)

run:
	ENV=dev go run $(MAIN_PATH)

test:
	go test ./... -v -count=1

test-cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run ./...

fmt:
	gofmt -w .
	goimports -w .

vet:
	go vet ./...

clean:
	rm -rf $(BUILD_DIR) coverage.out coverage.html

docker:
	docker build -t $(APP_NAME):latest -f Dockerfile.multistage .

docker-compose:
	docker compose up -d

docker-compose-down:
	docker compose down

swagger:
	swag init -g cmd/main.go -o docs/swagger
