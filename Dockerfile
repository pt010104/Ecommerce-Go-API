FROM golang:1.22.7-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/api/main.go
RUN go build -o api ./cmd/api

FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata && update-ca-certificates
ENV TZ=Asia/Ho_Chi_Minh
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
WORKDIR /app
COPY --from=builder /app/api .
EXPOSE 8080
ENTRYPOINT ["./api"]
