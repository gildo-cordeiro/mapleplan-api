# Dockerfile
FROM golang:1.25.4-alpine AS builder

# instalar certificados e git (necessário para baixar módulos de VCS)
RUN apk update && apk add --no-cache ca-certificates git

WORKDIR /app

ENV GOPROXY=https://proxy.golang.org,direct

# copiar dependências e baixar para aproveitar cache
COPY go.mod go.sum ./
RUN go mod download

# copiar código e compilar
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/main.go

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/main /main
EXPOSE 8080
ENTRYPOINT ["/main"]
