# Quick Start Guide - Auth Simple Ecommerce

Panduan cepat untuk menjalankan aplikasi authentication microservice.

## Prerequisites

Pastikan Anda sudah menginstall:
- Go 1.25+
- PostgreSQL 12+
- Make (optional, tapi recommended)

## Langkah-langkah Setup

### 1. Setup Database

```bash
# Buat database PostgreSQL
createdb auth_ecommerce

# Atau via psql
psql -U postgres
CREATE DATABASE auth_ecommerce;
\q
```

### 2. Setup Environment Variables

```bash
# Copy file .env.example
cp .env.example .env

# Edit .env sesuai konfigurasi Anda
# Minimal yang perlu diubah:
# - DB_PASSWORD=your_postgres_password
# - JWT_SECRET=your-random-secret-key
# - JWT_REFRESH_SECRET=your-random-refresh-secret-key
```

Contoh isi `.env`:
```env
PORT=9000
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=auth_ecommerce
DB_SSLMODE=disable
JWT_SECRET=my-super-secret-key-12345
JWT_REFRESH_SECRET=my-refresh-secret-key-67890
```

### 3. Install Dependencies

```bash
go mod download
```

### 4. Generate Swagger Docs (Optional)

```bash
make docs
# atau
swag init
```

### 5. Run Migrations (Optional)

GORM akan otomatis membuat tabel, tapi untuk field tambahan:

```bash
make migrate
# atau manual
psql -d auth_ecommerce -f migrations/001_add_auth_fields.sql
```

### 6. Run Application

```bash
# Menggunakan Make
make run

# Atau langsung dengan Go
go run main.go

# Atau build dulu
make build
./bin/app
```

Server akan berjalan di `http://localhost:9000`

## Testing API

### 1. Health Check

```bash
curl http://localhost:9000/api/
```

### 2. Register User

```bash
curl -X POST http://localhost:9000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "phone_number": "081234567890",
    "full_name": "John Doe",
    "password": "password123",
    "address": "Jakarta, Indonesia",
    "dob": "1990-01-01"
  }'
```

Response:
```json
{
  "code": 201,
  "message": "Registration successful",
  "data": {
    "user": {
      "id": 1,
      "username": "johndoe",
      "email": "john@example.com",
      ...
    },
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_at": "2025-10-31T12:00:00Z"
  }
}
```

### 3. Login

```bash
curl -X POST http://localhost:9000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email_or_username": "john@example.com",
    "password": "password123"
  }'
```

### 4. Get Profile (Protected Endpoint)

```bash
# Ganti YOUR_ACCESS_TOKEN dengan token dari response register/login
curl http://localhost:9000/api/v1/auth/profile \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 5. Logout

```bash
curl -X POST http://localhost:9000/api/v1/auth/logout \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## Menggunakan Swagger UI

Buka browser dan akses:
```
http://localhost:9000/api/swagger/index.html
```

Anda bisa testing semua endpoint langsung dari Swagger UI.

## Menggunakan Postman

1. Import file `Auth-API.postman_collection.json` ke Postman
2. Update variable `base_url` jika diperlukan (default: `http://localhost:9000/api`)
3. Jalankan request secara berurutan (Register -> Login -> endpoints lainnya)
4. Access token akan otomatis tersimpan di collection variables

## Troubleshooting

### Error: "failed to connect to database"

- Pastikan PostgreSQL sudah running
- Cek kredensial database di file `.env`
- Test koneksi manual: `psql -U postgres -d auth_ecommerce`

### Error: "JWT_SECRET not configured"

- Pastikan file `.env` sudah ada
- Pastikan variabel `JWT_SECRET` sudah diset
- Restart aplikasi setelah mengubah `.env`

### Error: port already in use

- Ganti PORT di `.env` ke port lain (misal 8080)
- Atau kill process yang menggunakan port 9000:
  ```bash
  lsof -ti:9000 | xargs kill -9
  ```

### Database migration errors

- Pastikan user PostgreSQL punya privilege yang cukup
- Jalankan migration manual jika perlu
- Cek log error untuk detail lebih lanjut

## Development Tips

### Hot Reload (menggunakan Air)

```bash
# Install air
go install github.com/air-verse/air@latest

# Run dengan hot reload
make dev
# atau
air
```

### Run dengan Docker

```bash
# Start services (app + postgres)
make docker-up

# View logs
make docker-logs

# Stop services
make docker-down
```

### Format Code

```bash
make lint
```

### Run Tests

```bash
make test

# Dengan coverage
make test-coverage
```

## Struktur Response

### Success Response
```json
{
  "code": 200,
  "message": "Success message",
  "data": { ... }
}
```

### Error Response
```json
{
  "code": 400,
  "message": "Error message",
  "data": {
    "details": "Additional error details"
  }
}
```

## Next Steps

1. Implementasi email service untuk forgot password
2. Tambahkan rate limiting
3. Implementasi refresh token rotation
4. Tambahkan logging yang lebih detail
5. Setup CI/CD pipeline
6. Tambahkan integration tests
7. Implementasi role-based permissions
8. Tambahkan user management endpoints

## Dokumentasi Lengkap

Lihat [README.md](README.md) untuk dokumentasi lengkap.

## Support

Jika ada pertanyaan atau masalah, silakan buat issue di repository.
