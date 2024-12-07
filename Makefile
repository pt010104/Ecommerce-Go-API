-include .env 
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

test: 
	@echo "Running tests..."
	@go test -mod=vendor -coverprofile=coverage.out -failfast -timeout 5m ./... 
	@grep -v 'mock_' coverage.out > c.out
	@go tool cover -func=c.out
	@echo "Coverage report generated in c.out"