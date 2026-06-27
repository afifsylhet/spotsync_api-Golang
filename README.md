# SpotSync API 🚗⚡

Smart Parking & EV Charging Reservation System

A centralized platform for busy airports and malls to manage parking zones, specifically handling the high-demand reservation of limited EV charging spots.

**GitHub**: [https://github.com/afifsylhet/spotsync_api-Golang](https://github.com/afifsylhet/spotsync_api-Golang)

---

## Features

- JWT authentication with role-based access (`driver` / `admin`)
- Manage parking zones (`general`, `ev_charging`, `covered`)
- Concurrency-safe reservations using DB transactions + row locks
- Clean modular architecture (Handler → Service → Repository)
- Global JSON error responses via `internal/httpresponse`
- Postman collection included for full API testing

---

## Tech Stack

| Technology | Purpose |
|------------|---------|
| **Go 1.22+** | Backend language |
| **Echo v4** | HTTP web framework |
| **PostgreSQL** | Relational database (NeonDB / Supabase) |
| **GORM** | ORM with auto-migrations |
| **JWT** | `golang-jwt/jwt/v5` token auth |
| **bcrypt** | Password hashing (cost 12) |
| **validator/v10** | Request validation |
| **Air** | Hot reload for local development |

---

## Architecture

```
HTTP Request
    ↓
Handler   (validates DTO, extracts JWT claims)
    ↓
Service   (business logic, rules)
    ↓
Repository (GORM database operations)
    ↓
PostgreSQL
```

### Project Structure

```text
spotsync-api/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── auth/                   # JWT, login, register
│   ├── user/                   # User entity, repository, service
│   ├── zone/                   # Parking zones module
│   ├── reservation/            # Reservations module (concurrency-safe)
│   ├── config/                 # Env loading & database connection
│   ├── httpresponse/           # Global JSON error/success helpers
│   └── server/                 # Echo setup, middleware, DI, routes
├── postman/
│   └── SpotSync-API.postman_collection.json
├── .env.example
├── go.mod
└── README.md
```

Each module follows the same layout:

- `entity.go` — GORM models
- `dto/request.go` & `dto/response.go` — API payloads
- `repository.go` — Database queries only
- `service.go` — Business logic
- `handler.go` — HTTP handlers
- `register.go` — Route registration

---

## Local Setup

1. **Clone the repo**

   ```bash
   git clone https://github.com/afifsylhet/spotsync_api-Golang.git
   cd spotsync-api
   ```

2. **Copy environment file**

   ```bash
   cp .env.example .env
   ```

   Fill in your PostgreSQL connection string and JWT secret.

3. **Install dependencies**

   ```bash
   go mod tidy
   ```

4. **Run the server**

   ```bash
   air
   ```

   Or without hot reload:

   ```bash
   go run ./cmd/main.go
   ```

   Server starts at `http://localhost:8080` by default.

---

## Required `.env` Variables

```env
DATABASE_URL=postgres://user:password@localhost:5432/spotsync?sslmode=disable
JWT_SECRET=your-secret-key
PORT=8080
```

| Variable | Description |
|----------|-------------|
| `DATABASE_URL` | PostgreSQL connection string |
| `JWT_SECRET` | Secret key for signing JWT tokens |
| `PORT` | HTTP port (default: `8080`) |

---

## API Endpoints

| Method | Endpoint | Access |
|--------|----------|--------|
| POST | `/api/v1/auth/register` | Public |
| POST | `/api/v1/auth/login` | Public |
| GET | `/api/v1/zones` | Public |
| GET | `/api/v1/zones/:id` | Public |
| POST | `/api/v1/zones` | Admin |
| PUT | `/api/v1/zones/:id` | Admin |
| DELETE | `/api/v1/zones/:id` | Admin |
| POST | `/api/v1/reservations` | Authenticated |
| GET | `/api/v1/reservations/my-reservations` | Authenticated |
| DELETE | `/api/v1/reservations/:id` | Authenticated |
| GET | `/api/v1/reservations` | Admin |

---

## User Roles

| Role | Permissions |
|------|-------------|
| **driver** | Register/login, view zones, reserve spots, view & cancel own reservations |
| **admin** | All driver permissions + CRUD zones, view all reservations |

---

## Response Format

**Success**

```json
{
  "success": true,
  "message": "Operation description",
  "data": {}
}
```

**Error**

```json
{
  "success": false,
  "message": "Error description",
  "errors": "Error details"
}
```

---

## Postman Collection

Import `postman/SpotSync-API.postman_collection.json` into Postman.

**Quick test flow:**

1. **Register Admin** → **Login Admin** (saves JWT automatically)
2. **Create Parking Zone** (saves `zoneId` automatically)
3. **Register Driver** → **Login Driver**
4. **Create Reservation** → **Get My Reservations** → **Cancel Reservation**

Collection variables: `baseUrl`, `token`, `zoneId`, `reservationId`

---

## Concurrency-Safe Reservations

Reservations use a GORM transaction with row-level locking (`FOR UPDATE`) on the parking zone record to prevent over-capacity bookings when multiple drivers reserve the last spot simultaneously.

---

## License

MIT
