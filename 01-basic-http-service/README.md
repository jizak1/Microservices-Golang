# 🚀 Basic HTTP Service

Service HTTP sederhana yang mendemonstrasikan konsep dasar microservice dengan Golang dan Gin framework.

## 📚 Apa yang Dipelajari

### 🎯 Konsep Dasar
- **REST API** dengan HTTP methods (GET, POST)
- **JSON** request dan response handling
- **Error handling** yang konsisten
- **Logging** dengan structured logging
- **Graceful shutdown** untuk production readiness

### 🛠️ Teknologi yang Digunakan
- **Gin** - HTTP web framework yang cepat
- **Logrus** - Structured logging
- **Standard Library** - HTTP server dan context

## 🏗️ Struktur Code

```
01-basic-http-service/
├── main.go              # Entry point aplikasi
├── README.md           # Dokumentasi ini
└── go.mod              # Go module dependencies
```

## 🚀 Cara Menjalankan

### Prerequisites
- Go 1.21+
- Git

### Langkah-langkah

1. **Clone dan masuk ke direktori**
```bash
cd 01-basic-http-service
```

2. **Initialize Go module**
```bash
go mod init basic-http-service
go mod tidy
```

3. **Install dependencies**
```bash
go get github.com/gin-gonic/gin
go get github.com/sirupsen/logrus
```

4. **Jalankan service**
```bash
go run main.go
```

Service akan berjalan di `http://localhost:8080`

## 📖 API Endpoints

### 🔍 Health Check
```http
GET /api/v1/health
```

**Response:**
```json
{
  "status": "healthy",
  "service": "basic-http-service",
  "timestamp": "2024-01-15T10:30:00Z",
  "version": "1.0.0"
}
```

### 📦 Get All Products
```http
GET /api/v1/products
```

**Response:**
```json
{
  "success": true,
  "message": "Products retrieved successfully",
  "data": [
    {
      "id": 1,
      "name": "Laptop Gaming",
      "description": "Laptop gaming dengan spesifikasi tinggi",
      "price": 15000000,
      "category": "Electronics",
      "created_at": "2024-01-15T10:30:00Z"
    }
  ],
  "count": 1
}
```

### 🔍 Get Product by ID
```http
GET /api/v1/products/{id}
```

**Response (Success):**
```json
{
  "success": true,
  "message": "Product retrieved successfully",
  "data": {
    "id": 1,
    "name": "Laptop Gaming",
    "description": "Laptop gaming dengan spesifikasi tinggi",
    "price": 15000000,
    "category": "Electronics",
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

**Response (Not Found):**
```json
{
  "success": false,
  "error": "Product not found",
  "message": "product with ID 999 not found"
}
```

### ➕ Create New Product
```http
POST /api/v1/products
Content-Type: application/json

{
  "name": "Smartphone Baru",
  "description": "Smartphone dengan fitur terbaru",
  "price": 5000000,
  "category": "Electronics"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Product created successfully",
  "data": {
    "id": 4,
    "name": "Smartphone Baru",
    "description": "Smartphone dengan fitur terbaru",
    "price": 5000000,
    "category": "Electronics",
    "created_at": "2024-01-15T10:35:00Z"
  }
}
```

## 🧪 Testing dengan cURL

### Test Health Check
```bash
curl http://localhost:8080/api/v1/health
```

### Test Get All Products
```bash
curl http://localhost:8080/api/v1/products
```

### Test Get Product by ID
```bash
curl http://localhost:8080/api/v1/products/1
```

### Test Create Product
```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Headphone Wireless",
    "description": "Headphone wireless dengan noise cancelling",
    "price": 1500000,
    "category": "Electronics"
  }'
```

## 🎨 Fitur Keren yang Diimplementasi

### 🔥 Clean Code Patterns
- **Separation of Concerns** - Handler, Service, dan Model terpisah
- **Dependency Injection** - Service di-inject ke Handler
- **Error Handling** - Consistent error response format
- **Input Validation** - Validasi input dengan error message yang jelas

### ⚡ Production Ready Features
- **Structured Logging** - JSON format untuk easy parsing
- **Graceful Shutdown** - Handle SIGINT/SIGTERM dengan proper cleanup
- **Health Check** - Endpoint untuk monitoring
- **HTTP Status Codes** - Proper HTTP status untuk setiap response

### 🛡️ Best Practices
- **JSON Response Format** - Consistent response structure
- **Error Messages** - User-friendly error messages
- **Request Validation** - Input validation dengan clear feedback
- **Logging Context** - Log dengan context information

## 🚀 Next Steps

Setelah memahami basic HTTP service ini, Anda bisa lanjut ke:

1. **02-user-management-service** - CRUD dengan database
2. **03-product-catalog-service** - Caching dan performance
3. **06-api-gateway** - Centralized routing

## 🤔 Pertanyaan Umum

**Q: Mengapa menggunakan Gin framework?**
A: Gin adalah framework yang cepat, ringan, dan mudah dipelajari. Perfect untuk microservices.

**Q: Bagaimana cara menambah middleware?**
A: Gunakan `router.Use(middleware)` sebelum mendefinisikan routes.

**Q: Bagaimana cara menambah validation yang lebih kompleks?**
A: Bisa menggunakan library seperti `go-playground/validator` atau custom validation functions.

---

**Happy Coding! 🎉**
