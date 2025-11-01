# Auth Simple Ecommerce - Authentication Service

Microservice authentication menggunakan Go, Echo Framework, GORM, PostgreSQL, dan JWT untuk sistem e-commerce.

## Features

✅ **Authentication & Authorization**
- User Registration dengan validasi
- Login dengan email/username
- JWT Token (Access Token & Refresh Token)
- Logout
- Profile Management

✅ **Password Management**
- Change Password (untuk user yang sudah login)
- Forgot Password (request reset token)
- Reset Password (dengan token)

✅ **Security**
- Password hashing dengan bcrypt
- JWT token validation
- Session management
- Role-based access control (RBAC)

✅ **Clean Architecture**
- Separation of concerns
- Dependency injection
- Interface-based design
- Standardized error handling

## Tech Stack

- **Language**: Go 1.25+
- **Framework**: Echo v4
- **ORM**: GORM
- **Database**: PostgreSQL
- **Authentication**: JWT (golang-jwt/jwt)
- **Validation**: go-playground/validator
- **API Documentation**: Swagger/OpenAPI

## Project Structure

```
.
├── cmd/                    # Application entry point
│   └── http.go            # HTTP server setup & routing
├── constants/             # Application constants
├── docs/                  # Swagger documentation
├── helpers/               # Helper functions
│   ├── config.go         # Configuration management
│   ├── errors.go         # Error handling
│   ├── jwt.go            # JWT utilities
│   ├── password.go       # Password hashing
│   ├── postgres.go       # Database connection
│   └── response.go       # HTTP response helpers
├── internal/
│   ├── api/              # HTTP handlers
│   │   ├── auth.go       # Auth endpoints
│   │   └── healthcheck.go
│   ├── interfaces/       # Interface definitions
│   │   └── IAuth.go
│   ├── middleware/       # HTTP middlewares
│   │   ├── auth.go       # JWT middleware
│   │   └── error.go      # Error handler middleware
│   ├── models/           # Domain models
│   │   ├── user.go
│   │   └── dto/          # Data Transfer Objects
│   │       └── auth.go
│   ├── repository/       # Data access layer
│   │   └── auth.go
│   └── services/         # Business logic layer
│       └── auth.go
├── migrations/           # Database migrations
├── .env.example          # Environment variables example
├── go.mod
└── main.go
```

## Getting Started

### Prerequisites

- Go 1.25 or higher
- PostgreSQL 12+
- Git

### Installation

1. **Clone the repository**
```bash
git clone <repository-url>
cd auth-simple-ecommerce
```

2. **Install dependencies**
```bash
go mod download
```

3. **Setup database**
```bash
# Create PostgreSQL database
createdb auth_ecommerce

# Run migrations (or let GORM auto-migrate)
psql -d auth_ecommerce -f migrations/001_add_auth_fields.sql
```

4. **Configure environment**
```bash
cp .env.example .env
# Edit .env with your configuration
```

Required environment variables:
```env
PORT=9000
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=auth_ecommerce
DB_SSLMODE=disable
JWT_SECRET=your-secret-key
JWT_REFRESH_SECRET=your-refresh-secret-key
```

5. **Run the application**
```bash
# Using Make (if Makefile exists)
make run

# Or directly with Go
go run main.go
```

The server will start at `http://localhost:9000`

## API Documentation

### Swagger/OpenAPI
Access interactive API documentation at:
```
http://localhost:9000/api/swagger/index.html
```

### API Endpoints

#### Public Endpoints (No Authentication Required)

**1. Register**
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "johndoe",
  "email": "john@example.com",
  "phone_number": "081234567890",
  "full_name": "John Doe",
  "password": "securePassword123",
  "address": "Jakarta, Indonesia",
  "dob": "1990-01-01"
}
```

**2. Login**
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email_or_username": "john@example.com",
  "password": "securePassword123"
}
```

Response:
```json
{
  "code": 200,
  "message": "Login successful",
  "data": {
    "user": {
      "id": 1,
      "username": "johndoe",
      "email": "john@example.com",
      "full_name": "John Doe",
      "role": "user"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_at": "2025-10-31T12:00:00Z"
  }
}
```

**3. Refresh Token**
```http
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**4. Forgot Password**
```http
POST /api/v1/auth/forgot-password
Content-Type: application/json

{
  "email": "john@example.com"
}
```

**5. Reset Password**
```http
POST /api/v1/auth/reset-password
Content-Type: application/json

{
  "token": "reset-token-from-email",
  "new_password": "newSecurePassword123"
}
```

#### Protected Endpoints (Authentication Required)

Add Authorization header: `Bearer <access_token>`

**6. Get Profile**
```http
GET /api/v1/auth/profile
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

**7. Change Password**
```http
POST /api/v1/auth/change-password
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
Content-Type: application/json

{
  "old_password": "currentPassword123",
  "new_password": "newSecurePassword123"
}
```

**8. Logout**
```http
POST /api/v1/auth/logout
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

**9. Health Check**
```http
GET /api/
```

## Error Handling

Standardized error response format:
```json
{
  "code": 400,
  "message": "Validation error",
  "data": {
    "details": "Username is required"
  }
}
```

Common HTTP Status Codes:
- `200` - Success
- `201` - Created
- `400` - Bad Request (validation errors)
- `401` - Unauthorized (invalid/expired token)
- `403` - Forbidden (insufficient permissions)
- `404` - Not Found
- `409` - Conflict (duplicate data)
- `500` - Internal Server Error

## Middleware

### JWT Middleware
Validates JWT token and extracts user information to context.

```go
// Usage in routes
authProtected := api.Group("/v1/auth")
authProtected.Use(appMiddleware.JWTMiddleware())
```

### Role Middleware
Restricts access based on user roles.

```go
// Usage example
admin := api.Group("/v1/admin")
admin.Use(appMiddleware.JWTMiddleware())
admin.Use(appMiddleware.RoleMiddleware("admin"))
```

### Error Handler Middleware
Standardizes all error responses.

## Security Best Practices

1. **Environment Variables**: Never commit `.env` file
2. **JWT Secret**: Use strong, random secrets in production
3. **Password Policy**: Minimum 6 characters (customize in validation)
4. **HTTPS**: Always use HTTPS in production
5. **Rate Limiting**: Consider adding rate limiting middleware
6. **CORS**: Configure CORS properly for your frontend

## Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test -run TestRegister ./internal/services
```

## Database Migrations

The application uses GORM's AutoMigrate feature, but you can also run manual migrations:

```bash
# Apply migrations
psql -d auth_ecommerce -f migrations/001_add_auth_fields.sql
```

## Development

### Code Structure Guidelines

1. **Repository Layer**: Database operations only
2. **Service Layer**: Business logic
3. **Handler/API Layer**: HTTP request/response handling
4. **Middleware**: Cross-cutting concerns (auth, logging, etc.)

### Adding New Features

1. Define interface in `internal/interfaces/`
2. Implement repository in `internal/repository/`
3. Implement service in `internal/services/`
4. Create handler in `internal/api/`
5. Register routes in `cmd/http.go`
6. Update Swagger documentation

## Deployment

### Using Docker

```bash
# Build image
docker build -t auth-service .

# Run container
docker run -p 9000:9000 --env-file .env auth-service
```

### Using Docker Compose

```bash
docker-compose -f docker-compose-dev.yaml up
```

## Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## License

[Add your license here]

## Contact

[Add your contact information]
