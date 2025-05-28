package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// AppConfig menyimpan konfigurasi aplikasi yang mudah dibaca
type AppConfig struct {
	// Server settings
	ServerPort     string
	ServerHost     string
	ServerTimeout  time.Duration
	ShutdownGracePeriod time.Duration
	
	// Database settings
	DatabaseURL      string
	DatabaseHost     string
	DatabasePort     string
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string
	DatabaseSSLMode  string
	MaxOpenConns     int
	MaxIdleConns     int
	ConnMaxLifetime  time.Duration
	
	// Redis settings
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int
	
	// JWT settings
	JWTSecret     string
	JWTExpiration time.Duration
	
	// External services
	UserServiceURL    string
	ProductServiceURL string
	OrderServiceURL   string
	NotificationServiceURL string
	
	// Message Queue settings
	RabbitMQURL      string
	KafkaBrokers     []string
	
	// Monitoring
	PrometheusPort   string
	JaegerEndpoint   string
	LogLevel         string
	
	// Environment
	Environment string
	Debug       bool
}

// LoadConfig memuat konfigurasi dari environment variables dan .env file
func LoadConfig() (*AppConfig, error) {
	// Load .env file jika ada (untuk development)
	_ = godotenv.Load()
	
	config := &AppConfig{
		// Server defaults
		ServerPort:     getEnvOrDefault("SERVER_PORT", "8080"),
		ServerHost:     getEnvOrDefault("SERVER_HOST", "localhost"),
		ServerTimeout:  getDurationOrDefault("SERVER_TIMEOUT", "30s"),
		ShutdownGracePeriod: getDurationOrDefault("SHUTDOWN_GRACE_PERIOD", "10s"),
		
		// Database defaults
		DatabaseHost:     getEnvOrDefault("DB_HOST", "localhost"),
		DatabasePort:     getEnvOrDefault("DB_PORT", "5432"),
		DatabaseUser:     getEnvOrDefault("DB_USER", "postgres"),
		DatabasePassword: getEnvOrDefault("DB_PASSWORD", "password"),
		DatabaseName:     getEnvOrDefault("DB_NAME", "microservices_db"),
		DatabaseSSLMode:  getEnvOrDefault("DB_SSL_MODE", "disable"),
		MaxOpenConns:     getIntOrDefault("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns:     getIntOrDefault("DB_MAX_IDLE_CONNS", 5),
		ConnMaxLifetime:  getDurationOrDefault("DB_CONN_MAX_LIFETIME", "5m"),
		
		// Redis defaults
		RedisHost:     getEnvOrDefault("REDIS_HOST", "localhost"),
		RedisPort:     getEnvOrDefault("REDIS_PORT", "6379"),
		RedisPassword: getEnvOrDefault("REDIS_PASSWORD", ""),
		RedisDB:       getIntOrDefault("REDIS_DB", 0),
		
		// JWT defaults
		JWTSecret:     getEnvOrDefault("JWT_SECRET", "your-super-secret-key-change-in-production"),
		JWTExpiration: getDurationOrDefault("JWT_EXPIRATION", "24h"),
		
		// Service URLs
		UserServiceURL:    getEnvOrDefault("USER_SERVICE_URL", "http://localhost:8081"),
		ProductServiceURL: getEnvOrDefault("PRODUCT_SERVICE_URL", "http://localhost:8082"),
		OrderServiceURL:   getEnvOrDefault("ORDER_SERVICE_URL", "http://localhost:8083"),
		NotificationServiceURL: getEnvOrDefault("NOTIFICATION_SERVICE_URL", "http://localhost:8084"),
		
		// Message Queue
		RabbitMQURL:  getEnvOrDefault("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		KafkaBrokers: getStringSliceOrDefault("KAFKA_BROKERS", []string{"localhost:9092"}),
		
		// Monitoring
		PrometheusPort: getEnvOrDefault("PROMETHEUS_PORT", "9090"),
		JaegerEndpoint: getEnvOrDefault("JAEGER_ENDPOINT", "http://localhost:14268/api/traces"),
		LogLevel:       getEnvOrDefault("LOG_LEVEL", "info"),
		
		// Environment
		Environment: getEnvOrDefault("ENVIRONMENT", "development"),
		Debug:       getBoolOrDefault("DEBUG", true),
	}
	
	// Build database URL
	config.DatabaseURL = fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.DatabaseUser,
		config.DatabasePassword,
		config.DatabaseHost,
		config.DatabasePort,
		config.DatabaseName,
		config.DatabaseSSLMode,
	)
	
	return config, nil
}

// Helper functions untuk parsing environment variables
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getBoolOrDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getDurationOrDefault(key string, defaultValue string) time.Duration {
	value := getEnvOrDefault(key, defaultValue)
	if duration, err := time.ParseDuration(value); err == nil {
		return duration
	}
	// Fallback ke default jika parsing gagal
	if duration, err := time.ParseDuration(defaultValue); err == nil {
		return duration
	}
	return time.Minute // Ultimate fallback
}

func getStringSliceOrDefault(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		// Simple split by comma - bisa diperbaiki dengan parser yang lebih sophisticated
		return []string{value}
	}
	return defaultValue
}

// IsProduction mengecek apakah aplikasi berjalan di production
func (c *AppConfig) IsProduction() bool {
	return c.Environment == "production"
}

// IsDevelopment mengecek apakah aplikasi berjalan di development
func (c *AppConfig) IsDevelopment() bool {
	return c.Environment == "development"
}

// GetServerAddress mengembalikan alamat lengkap server
func (c *AppConfig) GetServerAddress() string {
	return fmt.Sprintf("%s:%s", c.ServerHost, c.ServerPort)
}

// GetRedisAddress mengembalikan alamat lengkap Redis
func (c *AppConfig) GetRedisAddress() string {
	return fmt.Sprintf("%s:%s", c.RedisHost, c.RedisPort)
}
