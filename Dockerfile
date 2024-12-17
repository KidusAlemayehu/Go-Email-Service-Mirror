FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o api ./cmd/api
RUN go build -o worker ./cmd/worker

EXPOSE 8080

CMD ["./api"]
