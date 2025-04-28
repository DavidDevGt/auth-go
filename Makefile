APP_NAME=auth-go
VERSION ?= $(shell git describe --tags --always)
PORT ?= 8080

.PHONY: all run build tidy fmt vet lint test migrate help

all: build

run:
	PORT=$(PORT) go run main.go

build:
	mkdir -p bin
	go build -o bin/$(APP_NAME) main.go

tidy:
	go mod tidy

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	golangci-lint run

test:
	go test ./... -v

migrate:
	#goose -dir db/migrations up

help:
	@echo "Targets disponibles:"
	@echo "  make run       Ejecuta la app en dev"
	@echo "  make build     Compila el binario"
	@echo "  make fmt       Formatea el código"
	@echo "  make vet       Corre go vet"
	@echo "  make lint      Ejecuta golangci-lint"
	@echo "  make test      Ejecuta tests"
	@echo "  make migrate   Aplica migraciones"
	@echo "  make help      Muestra esta ayuda"