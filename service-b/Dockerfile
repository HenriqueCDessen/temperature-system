FROM golang:1.24.5-alpine AS builder
WORKDIR /app

# Usa módulo único da raiz
COPY go.mod go.sum ./
RUN go mod download

# Copia apenas o código do service-b
COPY service-b/ ./service-b/

# Compila o binário do Service B a partir da pasta do serviço
RUN go build -o /service-b ./service-b

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /service-b /service-b
EXPOSE 8081
ENTRYPOINT ["/service-b"]
