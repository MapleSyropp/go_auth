all: build

build:
	@echo "Building..."
	@go build -o main cmd/api/main.go

run:
	@templ generate
	@go run cmd/api/main.go
