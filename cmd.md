# ERP Internal Backend - CLI Cheat Sheet

## 1. Initial Setup

### Clone & Configure Environment

```powershell
# Copy environment file
Copy-Item .env.example .env

# (Optional) On bash/cmd:
# cp .env.example .env
```

### Install Go Dependencies

```powershell
go mod tidy
go mod download
```

### Install Database Migration CLI (`golang-migrate`)

```powershell
# Install via Go (with MySQL driver support)
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Verify installation
migrate -version
```

---

## 2. Running Application

### Start Development Server

```powershell
go run cmd/server/main.go
```

### Build Production Binary

```powershell
# Windows
go build -o bin/server.exe cmd/server/main.go

# Linux / Mac
go build -o bin/server cmd/server/main.go
```

---

## 3. Database Migrations (`golang-migrate`)

> **Note**: Replace `root:root@tcp(localhost:3306)/go-basic` with your MySQL connection credentials if different.

### Create New Migration File

```powershell
migrate create -ext sql -dir migration/mysql -seq <migration_name>
```

_Example:_

```powershell
migrate create -ext sql -dir migration/mysql -seq create_users_table
```

### Run Migrations (Up)

```powershell
# Apply all pending migrations
migrate -path migration/mysql -database "mysql://root:root@tcp(localhost:3306)/go-basic" -verbose up

# Apply N specific steps
migrate -path migration/mysql -database "mysql://root:root@tcp(localhost:3306)/go-basic" -verbose up 1
```

### Rollback Migrations (Down)

```powershell
# Rollback 1 migration step
migrate -path migration/mysql -database "mysql://root:root@tcp(localhost:3306)/go-basic" -verbose down 1

# Rollback all migrations (CAUTION)
migrate -path migration/mysql -database "mysql://root:root@tcp(localhost:3306)/go-basic" -verbose down
```

### Check Migration Status / Fix Dirty State

```powershell
# Check current migration version
migrate -path migration/mysql -database "mysql://root:root@tcp(localhost:3306)/go-basic" version

# Force migration version (fix dirty database flag after migration failure)
migrate -path migration/mysql -database "mysql://root:root@tcp(localhost:3306)/go-basic" force <version_number>
```

_Example to clear dirty state on version 1:_

```powershell
migrate -path migration/mysql -database "mysql://root:root@tcp(localhost:3306)/go-basic" force 1
```

---

## 4. Testing & Maintenance

### Run Tests

```powershell
# Run all unit and integration tests
go test -v ./...

# Run tests with coverage
go test -cover ./...
```

### Format Code

```powershell
# Format all Go files in workspace
go fmt ./...
```
