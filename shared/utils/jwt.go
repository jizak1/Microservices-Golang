package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWTClaims struktur claims untuk JWT token
type JWTClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// JWTManager untuk mengelola JWT tokens
type JWTManager struct {
	secretKey string
	issuer    string
}

// NewJWTManager membuat instance baru JWTManager
func NewJWTManager(secretKey, issuer string) *JWTManager {
	return &JWTManager{
		secretKey: secretKey,
		issuer:    issuer,
	}
}

// GenerateToken membuat JWT token baru
func (j *JWTManager) GenerateToken(userID, username, email, role string, expiration time.Duration) (string, error) {
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		Email:    email,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    j.issuer,
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken memvalidasi JWT token dan mengembalikan claims
func (j *JWTManager) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, fmt.Errorf("failed to parse claims")
	}

	return claims, nil
}

// RefreshToken membuat token baru dengan expiration yang diperpanjang
func (j *JWTManager) RefreshToken(tokenString string, newExpiration time.Duration) (string, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", fmt.Errorf("invalid token for refresh: %w", err)
	}

	// Buat token baru dengan expiration yang diperpanjang
	return j.GenerateToken(claims.UserID, claims.Username, claims.Email, claims.Role, newExpiration)
}

// ExtractUserID mengambil user ID dari token
func (j *JWTManager) ExtractUserID(tokenString string) (string, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}

// IsTokenExpired mengecek apakah token sudah expired
func (j *JWTManager) IsTokenExpired(tokenString string) bool {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return true
	}
	
	return claims.ExpiresAt.Time.Before(time.Now())
}

// GetTokenRemainingTime mengembalikan sisa waktu token
func (j *JWTManager) GetTokenRemainingTime(tokenString string) (time.Duration, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return 0, err
	}
	
	remaining := claims.ExpiresAt.Time.Sub(time.Now())
	if remaining < 0 {
		return 0, fmt.Errorf("token has expired")
	}
	
	return remaining, nil
}
