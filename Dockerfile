# Estágio 1: Compilação (Build)
FROM golang:1.23-alpine AS builder

# Instala certificados de segurança (importante para conexões externas com o Banco)
RUN apk update && apk add --no-cache ca-certificates

WORKDIR /app

# Copia os arquivos de dependências primeiro (aproveita o cache do Docker)
COPY go.mod go.sum ./
RUN go mod download

# Copia o resto do código
COPY . .

# Compila o binário estático (CGO_ENABLED=0 garante que rode em qualquer lugar)
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/main.go

# Estágio 2: Execução (Final)
FROM scratch

# Copia os certificados do estágio anterior
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copia o binário compilado
COPY --from=builder /app/main /main

# Porta que sua API escuta (ajuste se necessário)
EXPOSE 8080

# Comando para rodar
ENTRYPOINT ["/main"]