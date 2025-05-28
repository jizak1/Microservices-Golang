# 🚀 Microservices dengan Golang - Panduan Lengkap

Selamat datang di koleksi lengkap contoh kode microservices menggunakan Golang! Repository ini dirancang untuk membantu Anda mempelajari arsitektur microservices dari dasar hingga tingkat lanjut dengan pendekatan yang praktis dan mudah dipahami.

## 📚 Apa yang Akan Anda Pelajari

### 🎯 Konsep Dasar
- Arsitektur microservices vs monolith
- Service discovery dan communication
- API Gateway patterns
- Database per service
- Event-driven architecture

### 🛠️ Teknologi Modern
- **HTTP/REST APIs** dengan Gin framework
- **gRPC** untuk komunikasi antar service
- **Message Queues** dengan RabbitMQ dan Kafka
- **Database** PostgreSQL, MongoDB, Redis
- **Containerization** dengan Docker
- **Orchestration** dengan Docker Compose
- **Monitoring** dengan Prometheus dan Grafana
- **Distributed Tracing** dengan Jaeger

## 🏗️ Struktur Project

```
microservices-golang/
├── 01-basic-http-service/          # Service HTTP sederhana
├── 02-user-management-service/     # CRUD dengan database
├── 03-product-catalog-service/     # Service dengan caching
├── 04-order-processing-service/    # Event-driven service
├── 05-notification-service/        # Message queue consumer
├── 06-api-gateway/                 # Gateway dengan routing
├── 07-grpc-communication/          # gRPC antar service
├── 08-event-sourcing/              # Event sourcing pattern
├── 09-saga-pattern/                # Distributed transactions
├── 10-monitoring-observability/    # Metrics dan logging
├── docker-compose.yml              # Orchestration setup
├── shared/                         # Shared libraries
│   ├── config/                     # Configuration management
│   ├── database/                   # Database connections
│   ├── middleware/                 # Common middleware
│   └── utils/                      # Utility functions
└── scripts/                        # Automation scripts
```

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- PostgreSQL
- Redis
- RabbitMQ (opsional)

### Menjalankan Semua Services
```bash
# Clone repository
git clone <repository-url>
cd microservices-golang

# Start infrastructure
docker-compose up -d postgres redis rabbitmq

# Run all services
make run-all

# Atau jalankan service individual
cd 01-basic-http-service
go run main.go
```

## 📖 Panduan Belajar

### Level 1: Pemula
1. **Basic HTTP Service** - Memahami dasar REST API
2. **User Management** - CRUD operations dengan database
3. **Product Catalog** - Caching dan performance optimization

### Level 2: Menengah
4. **Order Processing** - Event-driven architecture
5. **Notification Service** - Message queues
6. **API Gateway** - Centralized routing dan authentication

### Level 3: Lanjutan
7. **gRPC Communication** - High-performance communication
8. **Event Sourcing** - Advanced data patterns
9. **Saga Pattern** - Distributed transactions
10. **Monitoring** - Observability dan debugging

## 🎨 Fitur Keren yang Akan Anda Pelajari

### 🔥 Modern Go Patterns
- **Clean Architecture** dengan dependency injection
- **Hexagonal Architecture** untuk testability
- **Repository Pattern** untuk data access
- **Factory Pattern** untuk service creation
- **Observer Pattern** untuk event handling

### ⚡ Performance Optimization
- Connection pooling
- Caching strategies (Redis, in-memory)
- Database indexing
- Async processing
- Load balancing

### 🛡️ Security Best Practices
- JWT authentication
- Rate limiting
- Input validation
- CORS handling
- Secret management

### 📊 Monitoring & Debugging
- Structured logging dengan logrus
- Metrics dengan Prometheus
- Health checks
- Distributed tracing
- Error handling patterns

## 🤝 Kontribusi

Kami sangat menghargai kontribusi Anda! Silakan:
1. Fork repository ini
2. Buat feature branch
3. Commit perubahan Anda
4. Push ke branch
5. Buat Pull Request

## 📝 Lisensi

MIT License - silakan gunakan untuk belajar dan proyek komersial.

## 🙋‍♂️ Bantuan

Jika Anda mengalami kesulitan atau memiliki pertanyaan:
- Buka issue di GitHub
- Lihat dokumentasi di setiap folder service
- Cek troubleshooting guide

---

**Happy Coding! 🎉**

Mari kita mulai perjalanan microservices Anda dengan Golang!
