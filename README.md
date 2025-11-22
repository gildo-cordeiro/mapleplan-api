# mapleplan-api — Arquitetura Hexagonal (Ports & Adapters)

Este repositório foi reorganizado para seguir a arquitetura Hexagonal (Ports & Adapters).

Estrutura principal (resumida):

```
cmd/meu-servico/main.go       # Ponto de entrada (usa a fábrica/registry para compor dependências)
internal/
  core/
    domain/                   # Entidades, agregados, value objects
    ports/
      repositories/           # Portas de saída (interfaces de repositório)
      services/               # Portas de entrada (interfaces de caso de uso)
  services/                   # Implementações dos casos de uso (dependem de ports)
  adapters/
    api/                      # Roteamento e registry de handlers
    repository/               # Implementações concretas dos repositórios (GORM)
    handlers/                 # Handlers HTTP (adaptadores de entrada)
    di/                       # Fábrica/registry para compor dependências
pkg/                          # Bibliotecas reutilizáveis
configs/                      # Configurações e exemplos (.env.example)
```

Como funciona
- As interfaces (ports) vivem em `internal/core/ports`.
- Implementações de infra (DB, HTTP, logger) vivem em `internal/adapters`.
- Os services (casos de uso) em `internal/services` dependem apenas das ports.
- A fábrica `internal/adapters/di/registry.go` compõe DB -> Repos -> Services -> Handlers
  e retorna um `api.HandlerRegistry` que é usado pelo roteador em `internal/adapters/api`.

Rodando localmente

1. Build:

```powershell
cd C:\Users\Gildo\mapleplan-api
go build ./...
```

2. Testes:

```powershell
go test ./...
```

3. Rodar o servidor (a partir da raiz):

```powershell
go run ./cmd/api
```

Próximos passos sugeridos
- Implementar os métodos pendentes do `TaskRepository` (adaptador GORM).
- Adicionar testes de integração para a camada HTTP.
- Adicionar um `Makefile` com comandos `make build`, `make run`, `make test`.

Se quiser, eu implemento os métodos do `TaskRepository` e adiciono uma fábrica para configuração (ex.: leitura de DB a partir de `configs`).
