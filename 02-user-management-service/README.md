# 👥 User Management Service

Service untuk mengelola user dengan operasi CRUD lengkap menggunakan PostgreSQL database. Mendemonstrasikan best practices untuk microservice dengan database integration.

## 📚 Apa yang Dipelajari

### 🎯 Konsep Database
- **PostgreSQL** integration dengan Go
- **CRUD Operations** (Create, Read, Update, Delete)
- **Database Migrations** otomatis
- **Connection Pooling** untuk performance
- **SQL Injection Prevention** dengan prepared statements
- **Database Indexing** untuk query optimization

### 🛠️ Advanced Patterns
- **Repository Pattern** untuk data access layer
- **Service Layer** untuk business logic
- **Input Validation** dengan Gin binding
- **Password Hashing** dengan bcrypt
- **Pagination** untuk large datasets
- **Error Handling** yang comprehensive

## 🏗️ Struktur Code

```
02-user-management-service/
├── main.go              # Complete service implementation
├── README.md           # Dokumentasi ini
├── go.mod              # Go module dependencies
└── docker-compose.yml  # PostgreSQL setup
```

## 🚀 Cara Menjalankan

### Prerequisites
- Go 1.21+
- PostgreSQL 13+
- Docker (opsional untuk database)

### Setup Database

#### Option 1: Menggunakan Docker
```bash
# Start PostgreSQL dengan Docker
docker run --name postgres-microservices \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=microservices_db \
  -p 5432:5432 \
  -d postgres:13
```

#### Option 2: PostgreSQL Lokal
```bash
# Install PostgreSQL dan buat database
createdb microservices_db
```

### Menjalankan Service

1. **Clone dan masuk ke direktori**
```bash
cd 02-user-management-service
```

2. **Initialize Go module**
```bash
go mod init user-management-service
go mod tidy
```

3. **Install dependencies**
```bash
go get github.com/gin-gonic/gin
go get github.com/jmoiron/sqlx
go get github.com/lib/pq
go get github.com/sirupsen/logrus
go get golang.org/x/crypto/bcrypt
```

4. **Jalankan service**
```bash
go run main.go
```

Service akan berjalan di `http://localhost:8081`

## 📖 API Endpoints

### 🔍 Health Check
```http
GET /api/v1/health
```

### 👤 Create User
```http
POST /api/v1/users
Content-Type: application/json

{
  "username": "johndoe",
  "email": "john@example.com",
  "full_name": "John Doe",
  "password": "securepassword123"
}
```

**Response:**
```json
{
  "success": true,
  "message": "User created successfully",
  "data": {
    "id": 1,
    "username": "johndoe",
    "email": "john@example.com",
    "full_name": "John Doe",
    "is_active": true,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

### 📋 Get All Users (dengan Pagination)
```http
GET /api/v1/users?page=1&limit=10
```

**Response:**
```json
{
  "success": true,
  "message": "Users retrieved successfully",
  "data": [
    {
      "id": 1,
      "username": "johndoe",
      "email": "john@example.com",
      "full_name": "John Doe",
      "is_active": true,
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 1,
    "total_pages": 1
  }
}
```

### 🔍 Get User by ID
```http
GET /api/v1/users/1
```

### ✏️ Update User
```http
PUT /api/v1/users/1
Content-Type: application/json

{
  "full_name": "John Smith",
  "is_active": false
}
```

### 🗑️ Delete User
```http
DELETE /api/v1/users/1
```

## 🧪 Testing dengan cURL

### Create User
```bash
curl -X POST http://localhost:8081/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "full_name": "Test User",
    "password": "password123"
  }'
```

### Get All Users
```bash
curl "http://localhost:8081/api/v1/users?page=1&limit=5"
```

### Get User by ID
```bash
curl http://localhost:8081/api/v1/users/1
```

### Update User
```bash
curl -X PUT http://localhost:8081/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Updated Name",
    "is_active": false
  }'
```

### Delete User
```bash
curl -X DELETE http://localhost:8081/api/v1/users/1
```

## 🎨 Fitur Keren yang Diimplementasi

### 🔥 Database Best Practices
- **Connection Pooling** - Optimal database connections
- **Prepared Statements** - SQL injection prevention
- **Database Indexes** - Fast query performance
- **Auto Migration** - Table creation otomatis
- **Transaction Support** - Data consistency

### ⚡ Advanced Features
- **Password Hashing** - bcrypt untuk security
- **Input Validation** - Comprehensive validation rules
- **Pagination** - Efficient large dataset handling
- **Unique Constraints** - Email dan username uniqueness
- **Soft Updates** - Partial update support

### 🛡️ Security & Validation
- **Email Validation** - Format email checking
- **Password Strength** - Minimum 6 characters
- **Username Rules** - 3-50 characters
- **SQL Injection Protection** - Parameterized queries
- **Error Sanitization** - Safe error messages

### 📊 Production Features
- **Structured Logging** - JSON format dengan context
- **Health Checks** - Database connectivity monitoring
- **Graceful Shutdown** - Proper resource cleanup
- **Error Handling** - Consistent error responses

## 🗃️ Database Schema

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes untuk performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_is_active ON users(is_active);
```

## 🚀 Next Steps

Setelah memahami user management service ini, Anda bisa lanjut ke:

1. **03-product-catalog-service** - Service dengan caching
2. **04-order-processing-service** - Event-driven architecture
3. **06-api-gateway** - Centralized authentication

## 🤔 Pertanyaan Umum

**Q: Bagaimana cara menambah field baru ke user?**
A: Update struct User, tambah field di database schema, dan update repository methods.

**Q: Bagaimana cara implement soft delete?**
A: Tambah field `deleted_at` dan update query untuk exclude deleted records.

**Q: Bagaimana cara menambah role-based access?**
A: Tambah field `role` dan implement middleware untuk authorization.

---

**Happy Coding! 🎉**
