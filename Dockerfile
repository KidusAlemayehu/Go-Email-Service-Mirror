FROM golang:1.23 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api cmd/api/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o worker cmd/worker/main.go

RUN chmod +x /app/api /app/worker

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY .env /app/.env 

COPY --from=builder /app/api /app/api
COPY --from=builder /app/worker /app/worker

EXPOSE ${BACKEND_PORT}
CMD ["./api"]
