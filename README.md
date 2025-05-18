# Weather Report Service - Deployment and Testing Guide

## Prerequisites
- Docker and Docker Compose installed
- Go 1.23 or later (for local development)
- PostgreSQL (automatically handled by Docker)
- Valid SSL certificates (for HTTPS)
- Weather API key from WeatherAPI.com

## Deployment Steps

### 1. Configuration Setup
1. Create a `config.json` file in the root directory.  
There is `example.config.json` that you can copy and feel respective fields. For example you should provide email and its corresponding [app-password](https://support.google.com/accounts/answer/185833?hl=en). Also you should have [weather-api](https://www.weatherapi.com/) token

2. Create `.env` file. Look at `.example.env` for reference.

3. You may place your SSL certificates in the `certs` directory:  
   - `certs/example.cert.pem`
   - `certs/example.key.pem`  

    However at present tls is not supported.

### 2. Deployment
```bash
# Build and start the services
docker-compose up --build -d
```

### 3. Migrations
```bash
# postgres forwards 5433:5432 port
# Default password defined in compose is "postgres"
psql -f ./migrations/up.sql -U postgres -h localhost -p 5433
```
`psql` tool should be installed.

