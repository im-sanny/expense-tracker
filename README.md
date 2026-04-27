# 📊 Expense Tracker API

A backend for expense tracking built with Go and PostgreSQL, focusing on clean architecture and secure authentication.

> 🔐 **Security First**: Implements protections against common vulnerabilities such as SQL injection and reduces risk of XSS/CSRF through secure cookie settings and token handling.

---

## ✨ Features

### 🔐 Authentication & Security

- **Secure Registration/Login**: Passwords hashed with `bcrypt` (cost 10).
- **Dual Token System**: Short-lived JWT Access Tokens (15 min) + Long-lived Refresh Tokens (7 days).
- **Cookie Security**: Cookies configured with HttpOnly, Secure, and SameSite attributes for safer token handling.
- **Token Revocation**: Logout invalidates refresh tokens server-side.
- **Generic Error Messages**: Prevents user enumeration attacks during login.

### 📝 Expense Management

- **Full CRUD**: Create, Read, Update (PUT/PATCH), Delete expenses.
- **Advanced Filtering**: Filter by date range, amount min/max, and search notes.
- **Pagination**: Efficient handling of large datasets.
- **User Isolation**: User isolation enforced at the application level (users can only access their own data).

### 🏗 Architecture & DevOps

- **Clean Architecture**: Separation of concerns (Handler → Service → Repository).
- **Database Migrations**: Version-controlled schema management with `golang-migrate`.
- **Structured Logging**: JSON logging with `log/slog` for observability.
- **Integration Testing**: Comprehensive test suite covering auth flow and CRUD operations.

---

## 🛠 Tech Stack

| Component      | Technology                       |
| -------------- | -------------------------------- |
| **Language**   | Go 1.25+                         |
| **Database**   | PostgreSQL 18 (with `uuid-ossp`) |
| **Auth**       | JWT (golang-jwt) + bcrypt        |
| **Migrations** | golang-migrate                   |
| **Testing**    | Bash + curl + jq                 |
| **Deployment** | Docker-ready, 12-Factor App      |

---

## 🚀 Quick Start

### Prerequisites

- Go 1.25+
- PostgreSQL 18+
- `golang-migrate` CLI (`go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`)
- `jq` for testing (`sudo apt install jq` or `brew install jq`)

### 1. Clone & Configure

```bash
git clone https://github.com/im-sanny/expense-tracker.git
cd expense-tracker
cp .env.example .env
```

Edit `.env`:

```ini
# Application
VERSION=1.0.0
SERVICE_NAME=EXPENSE_TRACKER
HTTP_PORT=3000
APP_ENV=local

# Database
DB_USER=postgres
DB_PASSWORD=your_password_here
DB_HOST=localhost
DB_PORT=5432
DB_NAME=etracker
DB_SSLMODE=disable

# Security (Generate with: openssl rand -base64 32)
JWT_SECRET_KEY=your-jwt-secret-min-32-chars
```

### 2. Setup Database

```bash
# Create database
createdb expense_tracker

# Run migrations
migrate -path migrations -database "postgres://postgres:yourpassword@localhost:5432/expense_tracker?sslmode=disable" up
```

### 3. Run Server

```bash
go run cmd/api/main.go
```

Server starts at `http://localhost:3000`

---

## 🧪 Testing

Run the automated integration test suite:

```bash
chmod +x test.sh
./test.sh
```

**Expected Output:**

```
🧪 Starting integration tests...
🔐 Testing Auth Flow...
✅ Register: 201
✅ Login: cookies saved
✅ Protected route: 200
📝 Testing CRUD Operations...
✅ Create: 201
✅ Get all: OK
✅ Update (PUT): OK
✅ Patch: OK
🎉 All tests passed!
```

---

## 📡 API Documentation

### Authentication

| Method | Endpoint         | Description             |
| ------ | ---------------- | ----------------------- |
| `POST` | `/auth/register` | Create new account      |
| `POST` | `/auth/login`    | Login + receive cookies |
| `POST` | `/auth/refresh`  | Get new access token    |
| `POST` | `/auth/logout`   | Revoke refresh token    |

### Expenses

| Method   | Endpoint      | Description                      |
| -------- | ------------- | -------------------------------- |
| `GET`    | `/track`      | List expenses (supports filters) |
| `POST`   | `/track`      | Create new expense               |
| `GET`    | `/track/{id}` | Get expense by ID                |
| `PUT`    | `/track/{id}` | Full update                      |
| `PATCH`  | `/track/{id}` | Partial update                   |
| `DELETE` | `/track/{id}` | Delete expense                   |

### Example Requests

**Register:**

```bash
curl -X POST http://localhost:3000/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"SecurePass123!"}'
```

**Create Expense:**

```bash
curl -X POST http://localhost:3000/track \
  -H "Content-Type: application/json" \
  -d '{"date":"2026-04-23T00:00:00Z","amount":890,"note":"Coffee"}' \
  -b cookies.txt
```

**Filter Expenses:**

```bash
# By date range
curl "http://localhost:3000/track?from=2026-01-01&to=2026-12-31" -b cookies.txt

# By amount
curl "http://localhost:3000/track?min=100&max=500" -b cookies.txt

# Search
curl "http://localhost:3000/track?q=coffee" -b cookies.txt
```

---

## 📂 Project Structure

```
/expense-tracker
├── cmd/expense-tracker             # Application entry point
├── internal/
│   ├── app/            # Routes
│   ├── config/         # Environment configuration
│   ├── database/       # DB connection & migrations
│   ├── handler/        # HTTP request/response logic
│   ├── service/        # Business logic
│   ├── repository/     # Data access layer
│   ├── model/          # Data structures
│   └── middlewares/    # Auth, CORS, Logging, etc.
├── migrations/         # SQL migration files
├── test.sh             # Integration test suite

```

---

## 🔒 Security Considerations

| Threat             | Mitigation                                   |
| ------------------ | -------------------------------------------- |
| **Password Leaks** | `bcrypt` hashing with salt                   |
| **XSS Attacks**    | `HttpOnly` cookies (JS cannot read tokens)   |
| **CSRF Attacks**   | `SameSite=Lax` cookie attribute              |
| **SQL Injection**  | Parameterized queries (`$1`, `$2`)           |
| **Brute Force**    | Generic error messages (prevent enumeration) |
| **Token Theft**    | Short-lived access tokens (15 min)           |

---

## 🚧 Deployment Note

This project is designed for **12-Factor App** deployment. To deploy to cloud platforms (Fly.io, Railway, Render):

1. Set `DATABASE_URL` environment variable.
2. Set `JWT_SECRET_KEY` (32+ random chars).
3. Set `ENV=production` (enables `Secure` cookie flag).
4. Run migrations on startup or via CI/CD.

> 📌 **Note**: For live demo, please contact the author or run locally. The codebase is fully tested and ready for deployment.

---
