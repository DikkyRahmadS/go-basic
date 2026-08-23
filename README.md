# ERP Internal Backend

A production-ready REST API built with **Go** using a **Modular Monolith** architecture.

---

## 🚀 Tech Stack

- **Language**: [Go 1.25+](https://golang.org/)
- **HTTP Framework**: [Gin Web Framework](https://github.com/gin-gonic/gin)
- **ORM**: [GORM](https://gorm.io/) (with MySQL driver)
- **Identifiers**: [google/uuid](https://github.com/google/uuid) for UUID v4 primary keys
- **Validation**: [go-playground/validator/v10](https://github.com/go-playground/validator)
- **Security**: [golang.org/x/crypto/bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) for password hashing
- **Configuration**: [godotenv](https://github.com/joho/godotenv)
- **Migrations**: [golang-migrate](https://github.com/golang-migrate/migrate)

---

## 📁 Project Structure

```
go-basic-backend/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── docs/
│   └── swagger.yml                 # OpenAPI 3.0 API specification
├── internal/
│   ├── config/                     # Core configs (App, DB, Gin, Validator, Env)
│   ├── middleware/                 # Shared HTTP middlewares (CORS)
│   │   └── cors.go
│   ├── modules/                    # Feature modules (Domain boundaries)
│   │   ├── health/                 # Health check module
│   │   │   ├── handler/
│   │   │   ├── models/
│   │   │   ├── router/
│   │   │   ├── service/
│   │   │   └── module.go
│   │   └── user/                   # User management module
│   │       ├── handler/            # HTTP request & response handlers
│   │       ├── models/             # Entities, DTOs, Converters
│   │       ├── repository/         # Database persistence layer
│   │       ├── router/             # Route definitions
│   │       ├── service/            # Business logic & domain rules
│   │       └── module.go           # Module constructor & DI wiring
│   └── pkg/                        # Shared reusable packages
│       ├── apperror/               # Centralized domain errors
│       ├── database/               # Context-aware transaction manager
│       ├── pagination/             # Generic pagination & GORM scope helper
│       └── response/               # Standardized JSON response envelope
├── migration/                      # Database migrations
│   └── mysql/                      # MySQL migration files (.up.sql & .down.sql)
├── .env.example                    # Environment template
├── ARCHITECTURE.md                 # Architectural guidelines & conventions
├── cmd.md                          # CLI & operations cheat sheet
├── go.mod
└── go.sum
```

---

## 🛠️ Getting Started

### 1. Prerequisites

- **Go** (version 1.25 or higher)
- **MySQL** (version 8.0+ or MariaDB 10.4+)
- **golang-migrate** CLI

### 2. Clone & Setup Environment

Copy the `.env.example` file and adjust your database credentials:

```bash
cp .env.example .env
```

Example `.env` configuration:

```env
# App Configuration
APP_PORT=8080
APP_ENV=development

# Database Configuration (MySQL)
DB_USER=root
DB_PASSWORD=root
DB_HOST=localhost
DB_PORT=3306
DB_NAME=erp_internal

# Database Pool Settings
DB_IDLE_CONNECTIONS=10
DB_MAX_CONNECTIONS=100
DB_MAX_LIFETIME_CONNECTIONS=3600

# CORS Configuration
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
```

### 3. Install Dependencies

```bash
go mod download
```

### 4. Run Database Migrations

```bash
migrate -path migration/mysql -database "mysql://root:root@tcp(localhost:3306)/erp_internal" -verbose up
```

### 5. Run the Application

```bash
go run cmd/server/main.go
```

The server will start listening at `http://localhost:8080`.

---

## 📖 API Documentation

The full OpenAPI / Swagger 3.0 specification is available at **[docs/swagger.yml](docs/swagger.yml)**.

### Available Endpoints:

- `GET /health` & `GET /api/health` — System and database health status
- `POST /api/users` — Create user
- `GET /api/users` — Get all users (supports optional search, filter, and pagination)
- `PUT /api/users/:id` — Update user
- `DELETE /api/users/:id` — Soft delete user

---

## 🏗️ Architecture & Best Practices

For in-depth architecture principles, cross-module dependency injection guidelines, and anti-patterns, refer to **[ARCHITECTURE.md](ARCHITECTURE.md)**.
For CLI command cheat sheets (migrations, build, testing), refer to **[cmd.md](cmd.md)**.
