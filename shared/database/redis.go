package database

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

// RedisClient wrapper untuk Redis connection yang mudah digunakan
type RedisClient struct {
	Client *redis.Client
	logger *logrus.Logger
}

// RedisConfig konfigurasi Redis yang user-friendly
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	Database int
}

// NewRedisConnection membuat koneksi baru ke Redis
func NewRedisConnection(config RedisConfig, logger *logrus.Logger) (*RedisClient, error) {
	address := fmt.Sprintf("%s:%s", config.Host, config.Port)
	
	logger.WithFields(logrus.Fields{
		"host":     config.Host,
		"port":     config.Port,
		"database": config.Database,
	}).Info("Connecting to Redis...")

	// Buat Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: config.Password,
		DB:       config.Database,
	})

	// Test koneksi
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		logger.WithError(err).Error("Failed to connect to Redis")
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.Info("Successfully connected to Redis")

	return &RedisClient{
		Client: client,
		logger: logger,
	}, nil
}

// Close menutup koneksi Redis
func (r *RedisClient) Close() error {
	if r.Client != nil {
		r.logger.Info("Closing Redis connection...")
		return r.Client.Close()
	}
	return nil
}

// HealthCheck mengecek kesehatan Redis connection
func (r *RedisClient) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	return r.Client.Ping(ctx).Err()
}

// SetWithExpiration menyimpan data dengan expiration time
func (r *RedisClient) SetWithExpiration(key string, value interface{}, expiration time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert value ke JSON jika bukan string
	var data string
	switch v := value.(type) {
	case string:
		data = v
	default:
		jsonData, err := json.Marshal(value)
		if err != nil {
			r.logger.WithError(err).WithField("key", key).Error("Failed to marshal value to JSON")
			return fmt.Errorf("failed to marshal value: %w", err)
		}
		data = string(jsonData)
	}

	if err := r.Client.Set(ctx, key, data, expiration).Err(); err != nil {
		r.logger.WithError(err).WithField("key", key).Error("Failed to set value in Redis")
		return fmt.Errorf("failed to set value: %w", err)
	}

	return nil
}

// Get mengambil data dari Redis
func (r *RedisClient) Get(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	value, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("key not found: %s", key)
		}
		r.logger.WithError(err).WithField("key", key).Error("Failed to get value from Redis")
		return "", fmt.Errorf("failed to get value: %w", err)
	}

	return value, nil
}

// GetAndUnmarshal mengambil data dan unmarshal ke struct
func (r *RedisClient) GetAndUnmarshal(key string, dest interface{}) error {
	value, err := r.Get(key)
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(value), dest); err != nil {
		r.logger.WithError(err).WithField("key", key).Error("Failed to unmarshal JSON")
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
}

// Delete menghapus key dari Redis
func (r *RedisClient) Delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := r.Client.Del(ctx, key).Err(); err != nil {
		r.logger.WithError(err).WithField("key", key).Error("Failed to delete key from Redis")
		return fmt.Errorf("failed to delete key: %w", err)
	}

	return nil
}

// Exists mengecek apakah key ada di Redis
func (r *RedisClient) Exists(key string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := r.Client.Exists(ctx, key).Result()
	if err != nil {
		r.logger.WithError(err).WithField("key", key).Error("Failed to check key existence")
		return false, fmt.Errorf("failed to check key existence: %w", err)
	}

	return count > 0, nil
}

// SetExpiration mengatur expiration untuk key yang sudah ada
func (r *RedisClient) SetExpiration(key string, expiration time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := r.Client.Expire(ctx, key, expiration).Err(); err != nil {
		r.logger.WithError(err).WithField("key", key).Error("Failed to set expiration")
		return fmt.Errorf("failed to set expiration: %w", err)
	}

	return nil
}

// GetTTL mendapatkan time-to-live untuk key
func (r *RedisClient) GetTTL(key string) (time.Duration, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ttl, err := r.Client.TTL(ctx, key).Result()
	if err != nil {
		r.logger.WithError(err).WithField("key", key).Error("Failed to get TTL")
		return 0, fmt.Errorf("failed to get TTL: %w", err)
	}

	return ttl, nil
}

// IncrementCounter increment counter dengan expiration
func (r *RedisClient) IncrementCounter(key string, expiration time.Duration) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Gunakan pipeline untuk atomic operation
	pipe := r.Client.Pipeline()
	incrCmd := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, expiration)
	
	if _, err := pipe.Exec(ctx); err != nil {
		r.logger.WithError(err).WithField("key", key).Error("Failed to increment counter")
		return 0, fmt.Errorf("failed to increment counter: %w", err)
	}

	return incrCmd.Val(), nil
}

// CacheWithCallback cache data dengan callback function jika cache miss
func (r *RedisClient) CacheWithCallback(key string, expiration time.Duration, callback func() (interface{}, error)) (interface{}, error) {
	// Coba ambil dari cache dulu
	var cachedData interface{}
	if err := r.GetAndUnmarshal(key, &cachedData); err == nil {
		r.logger.WithField("key", key).Debug("Cache hit")
		return cachedData, nil
	}

	// Cache miss, panggil callback
	r.logger.WithField("key", key).Debug("Cache miss, calling callback")
	data, err := callback()
	if err != nil {
		return nil, err
	}

	// Simpan ke cache
	if err := r.SetWithExpiration(key, data, expiration); err != nil {
		r.logger.WithError(err).WithField("key", key).Warn("Failed to cache data, but returning original data")
	}

	return data, nil
}
