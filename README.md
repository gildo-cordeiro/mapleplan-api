# MaplePlan API

Backend API for MaplePlan - a financial planning and goals management platform for couples planning to immigrate to Canada.

## 🏗️ Architecture: Simplified 3-Layer

A pragmatic architecture that balances separation of concerns with operational simplicity.

**3 Layers:**
1. **API Layer** → HTTP handlers, middlewares, routes
2. **Business Layer** → All pure business logic
3. **Data Layer** → Data access, storage, models

**Flow:** Request → Handlers (validate) → Services (logic) → Repositories (data) → Response

### 📁 Complete Structure

```
mapleplan-api/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   ├── middleware/
│   │   └── (routes and registry)
│   │
│   ├── business/
│   │   ├── (goal_service_impl.go)
│   │   ├── (user_service_impl.go)
│   │   └── (storage_service_impl.go)
│   │
│   ├── data/
│   │   ├── database/
│   │   ├── models/
│   │   │   ├── goal/
│   │   │   ├── user/
│   │   │   ├── couple/
│   │   │   ├── task/
│   │   │   ├── transaction/
│   │   │   └── province/
│   │   ├── repositories/
│   │   └── storage/
│   │
│   ├── dto/
│   │   ├── goal/
│   │   │   ├── request/
│   │   │   ├── response/
│   │   │   └── mapper/
│   │   └── user/
│   │       ├── request/
│   │       └── response/
│   │
│   ├── ports/
│   │   ├── repositories/
│   │   └── services/
│   │
│   └── bootstrap/
│       ├── (build_app.go)
│       ├── (build_storage.go)
│       └── (config.go)
│
├── pkg/
│   ├── jwt/
│   └── utils/
│
├── go.mod
├── go.sum
├── .env.example
├── docker-compose.yml
├── Dockerfile
└── README.md
```

## Getting Started

### Prerequisites
- Go 1.21+
- PostgreSQL
- AWS S3 (or MinIO)

### Setup
```bash
cp .env.example .env
go mod download
go build ./cmd/server/main.go
go run ./cmd/server/main.go
```

Server running at `http://localhost:8080`

## 🛠️ Stack

- **Framework**: Gorilla Mux
- **ORM**: GORM
- **Database**: PostgreSQL
- **Storage**: AWS S3
- **Auth**: JWT

## 📝 Conventions

- **Handlers**: `{Resource}Handler`
- **Services**: `{Resource}ServiceImpl`
- **Repositories**: `{Resource}RepositoryImpl`
- **Models**: `internal/data/models/{entity}`

## 🧪 Testing

```bash
go test ./...
go test -cover ./...
```

## 🤝 Contributing

```bash
git checkout -b feat/my-feature
git commit -m 'feat: description'
git push origin feat/my-feature
```