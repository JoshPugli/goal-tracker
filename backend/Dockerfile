# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Generate Swagger docs inside the build to ensure up-to-date API schema
RUN go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/server/main.go --output docs

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8000

CMD ["./main"]
