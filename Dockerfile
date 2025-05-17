FROM golang:1.23.9-alpine

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main ./cmd/app

EXPOSE 8080

CMD ["./main", "--host", "0.0.0.0", "--port", "8080", "--tls-certificate", "./certs/example.cert.pem", "--tls-key", "./certs/example.key.pem"]
