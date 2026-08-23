# Modular Monolith Architecture & Backend Guidelines

## 1. Directory Structure

This project follows an idiomatic Go **Modular Monolith (Package-by-Feature)** architecture:

```
go-basic-backend/
├── cmd/
│   └── server/
│       └── main.go                 # App entry point, bootstrap caller
├── docs/
│   └── swagger.yml                 # OpenAPI 3.0 API specifications
├── internal/
│   ├── config/                     # Application initializers & configurations
│   │   ├── app.go                  # Bootstrap & module wiring
│   │   ├── env.go                  # Env loader & typed getters
│   │   ├── gin.go                  # Gin router engine setup
│   │   ├── gorm.go                 # GORM database connection & connection pool
│   │   └── validator.go            # Go-playground validator setup
│   ├── middleware/                 # Shared HTTP middlewares
│   │   └── cors.go                 # Environment-aware CORS middleware
│   ├── modules/                    # Feature modules (Domain boundaries)
│   │   ├── health/                 # Health check module
│   │   │   ├── handler/
│   │   │   ├── models/
│   │   │   ├── router/
│   │   │   ├── service/
│   │   │   └── module.go
│   │   └── user/                   # User module
│   │       ├── handler/            # HTTP transport layer (Gin request/response bindings)
│   │       │   └── user_handler.go
│   │       ├── models/             # Entities, DTOs, & Converters
│   │       │   ├── user_entity.go
│   │       │   ├── user_dto.go
│   │       │   └── user_converter.go
│   │       ├── repository/         # Database persistence & queries
│   │       │   └── user_repository.go
│   │       ├── router/             # Module route definitions
│   │       │   └── user_router.go
│   │       ├── service/            # Core business logic & transactions
│   │       │   └── user_service.go
│   │       └── module.go           # Module constructor (NewModule) & DI wiring
│   └── pkg/                        # Reusable, single-responsibility internal packages
│       ├── apperror/               # Centralized domain error definitions & constructors
│       ├── database/               # Context-aware GORM transaction helper (WithTx, GetDB)
│       ├── pagination/             # Reusable pagination calculator & GORM scopes
│       └── response/               # Standardized JSON response envelopes & HandleError
├── migration/                      # Database schema migrations
│   └── mysql/                      # MySQL migration files (.up.sql & .down.sql)
├── .env.example
├── cmd.md                          # CLI cheat sheet
├── go.mod
└── go.sum
```

---

## 2. Layer Responsibilities

```
HTTP Request
    │
    ▼
[ Module.Router.RegisterRoutes ] (Gin Route Group)
    │
    ▼
[ Handler ]     --> Bind/validate JSON request payload (ShouldBindJSON / ShouldBindQuery)
    │               Parse HTTP headers / URL params / context
    │               Delegate error formatting to response.HandleError
    │               Return standard JSON response envelope (response.SuccessResponse)
    ▼
[ Service ]     --> Execute core business rules & domain logic
    │               Semantic validation (e.g. unique email checks)
    │               Password hashing (bcrypt)
    │               Manage DB transactions (database.WithTx)
    │               Return domain errors via apperror package
    ▼
[ Repository ]  --> Execute raw queries / GORM ORM operations
    │               Zero business logic, pure data access
    ▼
Database (MySQL)
```

---

## 3. Error Handling Architecture (`apperror` & `HandleError`)

To keep the service layer decoupled from HTTP transport frameworks (Gin), the application uses **Centralized Domain Errors**:

### 1. Service Layer (`apperror`)

Services return typed domain errors without knowing about HTTP requests or Gin contexts:

```go
// Domain conflict (409)
if existingUser != nil {
    return apperror.Conflict(fmt.Sprintf("email %s already exists", req.Email))
}

// Domain not found (404)
if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
    return apperror.NotFound("user not found")
}

// Low-level error wrapping (500)
if err != nil {
    return apperror.Internal(err)
}
```

Available constructors in `internal/pkg/apperror`:

- `apperror.BadRequest(msg)` (400)
- `apperror.Unauthorized(msg)` (401)
- `apperror.Forbidden(msg)` (403)
- `apperror.NotFound(msg)` (404)
- `apperror.Conflict(msg)` (409)
- `apperror.Internal(err)` (500)

### 2. Transport Layer (`response.HandleError`)

Handlers delegate all error responses to `response.HandleError(c, err)`:

```go
result, err := h.service.Create(c.Request.Context(), &req)
if err != nil {
    response.HandleError(c, err)
    return
}
```

`response.HandleError` automatically:

- Maps `*apperror.AppError` to its HTTP status code and message.
- Formats `validator.ValidationErrors` to `400 Bad Request` with field-level details.
- Defaults unknown/internal errors to `500 Internal Server Error` without exposing database internals.

---

## 4. Context-Aware Transactions (`database.WithTx`)

Atomic operations spanning multiple repository calls use `database.WithTx`:

```go
if err := database.WithTx(ctx, s.db, func(txCtx context.Context) error {
    existingUser, err := s.repository.FindByEmail(txCtx, user.Email)
    if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        return err
    }

    if existingUser != nil {
        return apperror.Conflict(fmt.Sprintf("email %s already exists", req.Email))
    }

    if err := s.repository.Create(txCtx, user); err != nil {
        return err
    }

    return nil
}); err != nil {
    return nil, err
}
```

Inside repositories, `database.GetDB(ctx, r.db)` automatically extracts the active transaction from context if one exists, or falls back to the default database connection.

---

## 5. Entity Standards: UUID & Soft Deletes

Entities must adhere to the following standards:

```go
type User struct {
    ID        string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
    Name      string         `gorm:"type:varchar(100);not null" json:"name"`
    Email     string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
    Password  string         `gorm:"type:varchar(255);not null" json:"-"`
    CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
    UpdatedAt time.Time      `gorm:"column:updated_at;autoCreateTime;autoUpdateTime" json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
    if u.ID == "" {
        u.ID = uuid.New().String()
    }
    return nil
}
```

- **UUID v4**: Auto-populated on `BeforeCreate` hook if empty.
- **Soft Deletes**: Uses `gorm.DeletedAt` for automatic soft deletion (`WHERE deleted_at IS NULL`).

---

## 6. Pagination & Filtering Standard

Reusable pagination and filtering logic is centralized in `internal/pkg/pagination`:

- **Repository**: Uses `pagination.Paginate(page, limit)` GORM Scope for SQL `LIMIT` and `OFFSET`.
- **Handler**: Uses `pagination.CalculateMeta(page, limit, total)` to build optional `response.Metadata`.

---

## 7. Cross-Module Communication Rules

### Rules Summary

| Allowed                                                       | Forbidden                                                 |
| ------------------------------------------------------------- | --------------------------------------------------------- |
| Module A calls Module B **Service Interface**                 | Module A imports Module B **Repository**                  |
| Module A receives Module B Service via constructor injection  | Module A queries Module B's database table directly       |
| Module A defines local interface for what it consumes         | Circular imports (Module A imports B, Module B imports A) |
| Store foreign entity **ID string** in entity structs          | Cross-import entity struct across module boundaries       |
| Async events (Message queue / Go channels) for loose coupling | Direct shared mutable memory state                        |

---

### Anti-Patterns to Avoid

- ❌ **Repository Leaking**: A service directly querying a foreign module's repository.
- ❌ **Cross-Domain Entity Preload**: Directly embedding another module's entity struct in models.
- ❌ **Circular Dependencies**: Module A importing B and Module B importing A.
