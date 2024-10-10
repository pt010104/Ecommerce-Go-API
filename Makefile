include .env
export
BINARY=engine
## Run the application
run-api:
	@echo "Running the application"
	@make swagger
	@go run cmd/api/main.go

swagger:
	@echo "Generating swagger"
	@swag init -g cmd/api/main.go
