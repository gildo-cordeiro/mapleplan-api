# AGENTS.md

**For: GitHub Copilot, Cursor, Claude, and AI code assistants**

## Quick Commands

```bash
go run cmd/server/main.go                    # Start dev server (http://localhost:8080)
go test ./...                                 # Run all tests
go mod tidy                                   # Clean dependencies
```

---

## Stack & Architecture

- **Language**: Go 1.21+
- **Database**: PostgreSQL + GORM (no joins in queries, use Preload)
- **Auth**: JWT (HS256, pkg/jwt)
- **Routing**: gorilla/mux
- **Storage**: AWS S3 (via AWS SDK v2)
- **Architecture**: Hexagonal (ports/interfaces + implementations)
- **Error Handling**: Sentinel errors in pkg/utils

**Layers** (request flow):
```
Handlers (net/http) → Services (business logic) → Repositories (GORM)
     ↓                      ↓                            ↓
api/handlers/          internal/business/      internal/data/repositories/
  + Middleware (auth)      internal/ports/      internal/data/models/
```

---

## Rules & Style

### 1. Handlers (api/handlers/)

**Pattern**: Plain struct with injected service interfaces.

```go
// GOOD - handler pattern
type GoalHandler struct {
    GoalService servicePort.GoalService
}

func (h *GoalHandler) CreateGoal(w http.ResponseWriter, r *http.Request) {
    userID := middleware.GetUserIDFromContext(r.Context())
    if userID == "" {
        http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
        return
    }
    
    var req dto.CreateGoalRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
        return
    }
    
    goal, err := h.GoalService.CreateGoal(r.Context(), userID, req)
    if err != nil {
        // Map errors to HTTP codes
        if errors.Is(err, utils.ErrAlreadyExists) {
            http.Error(w, `{"error":"already exists"}`, http.StatusConflict)
            return
        }
        http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(goal)
}
```

**Rules**:
- Extract `userID` from context: `middleware.GetUserIDFromContext(r.Context())`
- Always decode request body with error checking
- Map domain errors to HTTP status codes explicitly
- Use `http.StatusCode` constants (not magic numbers)
- Set `Content-Type: application/json` before writing response

### 2. Services (internal/business/)

**Pattern**: Implement ports (interfaces), orchestrate business logic.

```go
// GOOD - service pattern
type GoalService struct {
    goalRepo    ports.GoalRepository
    userRepo    ports.UserRepository
    txManager   ports.TransactionManager
}

func (s *GoalService) CreateGoal(ctx context.Context, userID string, req dto.CreateGoalRequest) (*dto.GoalResponse, error) {
    // 1. Validate input
    if req.Name == "" {
        return nil, errors.New("name is required")
    }
    
    // 2. Check permissions (user exists)
    user, err := s.userRepo.FindByID(ctx, userID)
    if err != nil {
        return nil, err
    }
    if user == nil {
        return nil, fmt.Errorf("user not found")
    }
    
    // 3. Transactional operation
    var goal *models.Goal
    err = s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
        goal = &models.Goal{
            UserID:   &userID,
            Name:     req.Name,
            Status:   models.GoalStatusActive,
        }
        return s.goalRepo.Save(txCtx, goal)
    })
    if err != nil {
        return nil, err
    }
    
    return dto.NewGoalResponse(goal), nil
}
```

**Rules**:
- Services operate on domain models (`models.*`), not DTOs
- Always validate input before repository calls
- Use transactions via `txManager.WithTransaction()` for multi-step operations
- Pass context through all calls: `WithTransaction(ctx, func(txCtx ...)`
- Map domain models to DTOs only in return statements

### 3. Repositories (internal/data/repositories/)

**Pattern**: GORM queries with transaction-aware DB selection.

```go
// GOOD - repository pattern
type GoalRepository struct {
    db *gorm.DB
}

// Helper: use transaction DB if available in context
func (r *GoalRepository) getDB(ctx context.Context) *gorm.DB {
    if tx, ok := ctx.Value("tx_key").(*gorm.DB); ok {
        return tx.WithContext(ctx)
    }
    return r.db.WithContext(ctx)
}

func (r *GoalRepository) FindByID(ctx context.Context, goalID string) (*models.Goal, error) {
    var goal models.Goal
    // ALWAYS use Preload for associations, never JOIN
    err := r.getDB(ctx).
        Preload("User").
        Preload("ImmigrationProfile").
        First(&goal, "id = ?", goalID).Error
    
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil // Consistent: nil, nil for not found
    }
    if err != nil {
        return nil, err
    }
    return &goal, nil
}

func (r *GoalRepository) Save(ctx context.Context, goal *models.Goal) error {
    return r.getDB(ctx).Save(goal).Error
}

func (r *GoalRepository) ListByUserID(ctx context.Context, userID string) ([]models.Goal, error) {
    var goals []models.Goal
    err := r.getDB(ctx).
        Where("user_id = ?", userID).
        Preload("User").
        Order("created_at DESC").
        Find(&goals).Error
    
    return goals, err
}
```

**Rules**:
- Always call `getDB(ctx)` for transaction-aware queries
- Use `Preload()` (not `Join()`) for associations
- Return `nil, nil` for not-found (don't wrap `gorm.ErrRecordNotFound`)
- Add `WHERE` clauses for scope (e.g., user_id filter)
- Handle Postgres errors explicitly: inspect `pq.Error` for constraint violations (code "23505" = duplicate key)

### 4. Models (internal/data/models/)

**Pattern**: Embed `models.Base`, use GORM tags, validation hooks.

```go
// GOOD - model pattern
type Goal struct {
    models.Base
    UserID                *string `gorm:"type:uuid;index" json:"user_id"`
    ImmigrationProfileID  *string `gorm:"type:uuid;index" json:"immigration_profile_id"`
    Name                  string  `gorm:"type:varchar(255);not null" json:"name"`
    Status                string  `gorm:"type:varchar(50);default:'active'" json:"status"`
    Progress              int     `gorm:"default:0" json:"progress"`
    DueDate               *time.Time `json:"due_date"`
    
    // Associations (not persisted, only for Preload)
    User                  *User              `gorm:"foreignKey:UserID" json:"user,omitempty"`
    ImmigrationProfile    *ImmigrationProfile `gorm:"foreignKey:ImmigrationProfileID" json:"profile,omitempty"`
}

// BeforeSave hook: validate domain constraints
func (g *Goal) BeforeSave(tx *gorm.DB) error {
    if g.UserID != nil && g.ImmigrationProfileID != nil {
        return errors.New("goal cannot belong to both user and profile")
    }
    if g.UserID == nil && g.ImmigrationProfileID == nil {
        return errors.New("goal must belong to user or profile")
    }
    return nil
}
```

**Rules**:
- Embed `models.Base` (provides UUID, timestamps, soft delete)
- Use GORM struct tags for column types, constraints, indexes
- Add associations as pointers (e.g., `*User`) with `gorm:"foreignKey"` tag
- Use `BeforeSave()` hooks for domain validations, not DB constraints
- JSON tags should match camelCase field names

### 5. Authentication & Context

**Pattern**: JWT in Authorization header or cookie, passed via context.

```go
// GOOD - middleware usage
func GetUserIDFromContext(ctx context.Context) string {
    userID, _ := ctx.Value(middleware.CtxUserIDKey).(string)
    return userID
}

// In handler:
userID := middleware.GetUserIDFromContext(r.Context())
if userID == "" {
    http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
    return
}

// In service, pass context down:
goals, err := s.goalRepo.ListByUserID(ctx, userID) // ctx already has userID
```

**Rules**:
- Never hardcode permissions; always extract userID from context
- Validate userID is non-empty before DB queries
- Pass context to all repository calls
- Use JWT_SECRET env var (not hardcoded)

### 6. Error Handling

**Pattern**: Return sentinel errors from repositories, map in handlers.

```go
// GOOD - error pattern
// In repository (data layer):
if err := r.db.Create(user).Error; err != nil {
    var pqErr *pq.Error
    if errors.As(err, &pqErr) && pqErr.Code == "23505" {
        return "", utils.ErrAlreadyExists // Sentinel error
    }
    return "", fmt.Errorf("create user: %w", err)
}

// In handler (api layer):
_, err := h.UserService.CreateUser(r.Context(), req)
if err != nil {
    if errors.Is(err, utils.ErrAlreadyExists) {
        http.Error(w, `{"error":"user already exists"}`, http.StatusConflict)
        return
    }
    http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
    return
}
```

**Rules**:
- Define sentinel errors in `pkg/utils/errors.go`
- Return `fmt.Errorf("context: %w", err)` for wrapping
- Use `errors.Is()` to check sentinel errors (not string matching)
- Map errors to HTTP codes only in handlers, never in services

---

## DO NOT

- ❌ Join tables in GORM queries. Use `Preload()` instead.
- ❌ Hardcode secrets or API keys. Use env vars via `config.LoadConfig()`.
- ❌ Write business logic in handlers. Move it to services.
- ❌ Return database models directly from handlers. Use DTOs + mappers.
- ❌ Ignore context cancellation. Pass `ctx` to all repository/service calls.
- ❌ Create ad-hoc error strings. Use sentinel errors from `pkg/utils`.
- ❌ Commit transactions in repositories. Let services/txManager handle it.
- ❌ Skip input validation. Check required fields before DB queries.
- ❌ Mix authentication checks with business logic. Do it at the handler layer.
- ❌ Add external dependencies without justification. Stick to: gorm, aws-sdk-go, gorilla/mux, jwt, logrus.

---

## Plan-Act-Reflect

When implementing a feature:

1. **Plan**: Outline the request flow:
   - Handler receives request
   - Service validates and orchestrates
   - Repository executes GORM query with transaction if needed
   - Handler maps response to DTO

2. **Act**: Write code in order:
   - Define domain model (if new)
   - Write repository method with `getDB(ctx)` + `Preload()`
   - Write service method with validation + transaction
   - Write handler with error mapping
   - Register route in `api/register_routes.go`

3. **Reflect**: Check:
   - [ ] All functions accept `context.Context`
   - [ ] userID extracted in handler and passed to service
   - [ ] Transactions used for multi-step operations
   - [ ] Errors mapped to correct HTTP codes
   - [ ] DTOs used in handlers, models in services
   - [ ] No `Join()` queries; only `Preload()`
   - [ ] Sentinel errors used, not error strings

---

## File Structure Commands

```bash
# New handler
touch internal/api/handlers/feature_handler.go

# New service
touch internal/business/feature_service_impl.go

# New repository
touch internal/data/repositories/feature_repository_impl.go

# New model
touch internal/data/models/feature/feature.go

# New port (interface)
touch internal/ports/services/feature_service_port.go
touch internal/ports/repositories/feature_repository_port.go
```

---

## Testing

- Tests should follow `_test.go` suffix.
- Use `TestXxx` naming for unit tests.
- Mock dependencies via struct composition and interface injection.
- Example test structure (handlers accept interfaces, easy to mock):

```go
type MockGoalService struct{}
func (m *MockGoalService) CreateGoal(ctx context.Context, userID string, req dto.CreateGoalRequest) (*dto.GoalResponse, error) {
    return &dto.GoalResponse{ID: "test-id", Name: req.Name}, nil
}

func TestCreateGoalHandler(t *testing.T) {
    handler := &GoalHandler{GoalService: &MockGoalService{}}
    // ... test logic
}
```

---

**Last Updated**: 2026  
**For**: Go 1.21+, PostgreSQL, GORM, JWT, AWS S3  
**Maintainer**: Team
