package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

// Logger middleware untuk logging request yang mudah dibaca
func Logger(logger *logrus.Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logger.WithFields(logrus.Fields{
			"client_ip":   param.ClientIP,
			"method":      param.Method,
			"path":        param.Path,
			"status_code": param.StatusCode,
			"latency":     param.Latency,
			"user_agent":  param.Request.UserAgent(),
		}).Info("HTTP Request")

		return ""
	})
}

// CORS middleware untuk handling Cross-Origin requests
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// RateLimiter middleware untuk membatasi request rate
func RateLimiter(requestsPerSecond int, burstSize int) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(requestsPerSecond), burstSize)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Rate limit exceeded",
				"message": fmt.Sprintf("Maximum %d requests per second allowed", requestsPerSecond),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// JWTAuth middleware untuk authentication dengan JWT
func JWTAuth(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Authorization header required",
				"message": "Please provide Authorization header with Bearer token",
			})
			c.Abort()
			return
		}

		// Extract token dari "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid authorization format",
				"message": "Authorization header must be in format: Bearer <token>",
			})
			c.Abort()
			return
		}

		// Parse dan validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid token",
				"message": "Token is invalid or expired",
			})
			c.Abort()
			return
		}

		// Extract claims dan simpan ke context
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("user_id", claims["user_id"])
			c.Set("username", claims["username"])
			c.Set("email", claims["email"])
		}

		c.Next()
	}
}

// RequestID middleware untuk menambahkan unique request ID
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next()
	}
}

// Timeout middleware untuk membatasi waktu request
func Timeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set timeout untuk context
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// Recovery middleware dengan logging yang informatif
func Recovery(logger *logrus.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logger.WithFields(logrus.Fields{
			"panic":      recovered,
			"path":       c.Request.URL.Path,
			"method":     c.Request.Method,
			"client_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		}).Error("Panic recovered")

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"message": "Something went wrong, please try again later",
		})
	})
}

// HealthCheck middleware untuk endpoint health check
func HealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/health" {
			c.JSON(http.StatusOK, gin.H{
				"status":    "healthy",
				"timestamp": time.Now().UTC(),
				"service":   "microservice",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// SecurityHeaders middleware untuk menambahkan security headers
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Next()
	}
}

// RequestSize middleware untuk membatasi ukuran request body
func RequestSize(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxSize {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error":   "Request too large",
				"message": fmt.Sprintf("Request body must be smaller than %d bytes", maxSize),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// Helper function untuk generate request ID
func generateRequestID() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), time.Now().Unix())
}

// GetUserFromContext helper untuk mengambil user info dari context
func GetUserFromContext(c *gin.Context) (userID string, username string, email string, exists bool) {
	userIDInterface, exists1 := c.Get("user_id")
	usernameInterface, exists2 := c.Get("username")
	emailInterface, exists3 := c.Get("email")

	if !exists1 || !exists2 || !exists3 {
		return "", "", "", false
	}

	userID, ok1 := userIDInterface.(string)
	username, ok2 := usernameInterface.(string)
	email, ok3 := emailInterface.(string)

	if !ok1 || !ok2 || !ok3 {
		return "", "", "", false
	}

	return userID, username, email, true
}
