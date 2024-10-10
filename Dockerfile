# syntax=docker/dockerfile:1

##
## STEP 1 - BUILD
##

# Use the official Golang image (Alpine-based) as the build environment
FROM golang:1.22.7-alpine AS builder

# Install necessary packages
RUN apk add --no-cache git

# Set the working directory inside the image
WORKDIR /app

# Copy go.mod and go.sum files from the build context to /app
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project to the working directory
COPY . .

# Install the swag CLI for generating Swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Generate the Swagger documentation
RUN swag init -g cmd/api/main.go

# Build the Go application binary
RUN go build -o api ./cmd/api

##
## STEP 2 - DEPLOY
##

# Use a minimal Alpine Linux image for the final deployment
FROM alpine:latest

# Install ca-certificates and tzdata (for timezone settings)
RUN apk add --no-cache ca-certificates tzdata && update-ca-certificates

# Set the timezone
ENV TZ=Asia/Ho_Chi_Minh
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Set the working directory inside the container
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/api .

# Expose the port that the application will listen on
# Note: Heroku assigns the port via the PORT environment variable,
# so the EXPOSE instruction can remain as is or be adjusted if desired.
EXPOSE 8080

# Run the Go application
ENTRYPOINT ["./api"]
