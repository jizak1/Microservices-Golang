package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// User model untuk database
type User struct {
	ID        int       `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	FullName  string    `json:"full_name" db:"full_name"`
	Password  string    `json:"-" db:"password_hash"` // Hidden dari JSON response
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CreateUserRequest untuk request body create user
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	FullName string `json:"full_name" binding:"required,min=2,max=100"`
	Password string `json:"password" binding:"required,min=6"`
}

// UpdateUserRequest untuk request body update user
type UpdateUserRequest struct {
	Username string `json:"username,omitempty" binding:"omitempty,min=3,max=50"`
	Email    string `json:"email,omitempty" binding:"omitempty,email"`
	FullName string `json:"full_name,omitempty" binding:"omitempty,min=2,max=100"`
	IsActive *bool  `json:"is_active,omitempty"`
}

// UserRepository untuk database operations
type UserRepository struct {
	db     *sqlx.DB
	logger *logrus.Logger
}

// NewUserRepository membuat instance baru UserRepository
func NewUserRepository(db *sqlx.DB, logger *logrus.Logger) *UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

// CreateUser menyimpan user baru ke database
func (ur *UserRepository) CreateUser(user *User) error {
	query := `
		INSERT INTO users (username, email, full_name, password_hash, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`

	err := ur.db.QueryRow(query, user.Username, user.Email, user.FullName,
		user.Password, user.IsActive, user.CreatedAt, user.UpdatedAt).Scan(&user.ID)

	if err != nil {
		ur.logger.WithError(err).Error("Failed to create user")
		return fmt.Errorf("failed to create user: %w", err)
	}

	ur.logger.WithField("user_id", user.ID).Info("User created successfully")
	return nil
}

// GetUserByID mengambil user berdasarkan ID
func (ur *UserRepository) GetUserByID(id int) (*User, error) {
	var user User
	query := `SELECT id, username, email, full_name, password_hash, is_active, created_at, updated_at
			  FROM users WHERE id = $1`

	err := ur.db.Get(&user, query, id)
	if err != nil {
		ur.logger.WithError(err).WithField("user_id", id).Error("Failed to get user")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetAllUsers mengambil semua users dengan pagination
func (ur *UserRepository) GetAllUsers(limit, offset int) ([]User, int, error) {
	var users []User
	var total int

	// Get total count
	countQuery := `SELECT COUNT(*) FROM users`
	err := ur.db.Get(&total, countQuery)
	if err != nil {
		ur.logger.WithError(err).Error("Failed to get user count")
		return nil, 0, fmt.Errorf("failed to get user count: %w", err)
	}

	// Get users with pagination
	query := `SELECT id, username, email, full_name, password_hash, is_active, created_at, updated_at
			  FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`

	err = ur.db.Select(&users, query, limit, offset)
	if err != nil {
		ur.logger.WithError(err).Error("Failed to get users")
		return nil, 0, fmt.Errorf("failed to get users: %w", err)
	}

	return users, total, nil
}

// UpdateUser mengupdate user di database
func (ur *UserRepository) UpdateUser(id int, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return fmt.Errorf("no fields to update")
	}

	// Build dynamic query
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	for field, value := range updates {
		setParts = append(setParts, fmt.Sprintf("%s = $%d", field, argIndex))
		args = append(args, value)
		argIndex++
	}

	// Add updated_at
	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	// Add WHERE clause
	args = append(args, id)

	query := fmt.Sprintf("UPDATE users SET %s WHERE id = $%d",
		joinStrings(setParts, ", "), argIndex)

	result, err := ur.db.Exec(query, args...)
	if err != nil {
		ur.logger.WithError(err).WithField("user_id", id).Error("Failed to update user")
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	ur.logger.WithField("user_id", id).Info("User updated successfully")
	return nil
}

// DeleteUser menghapus user dari database
func (ur *UserRepository) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := ur.db.Exec(query, id)
	if err != nil {
		ur.logger.WithError(err).WithField("user_id", id).Error("Failed to delete user")
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	ur.logger.WithField("user_id", id).Info("User deleted successfully")
	return nil
}

// CheckEmailExists mengecek apakah email sudah ada
func (ur *UserRepository) CheckEmailExists(email string, excludeID int) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE email = $1 AND id != $2`

	err := ur.db.Get(&count, query, email, excludeID)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// CheckUsernameExists mengecek apakah username sudah ada
func (ur *UserRepository) CheckUsernameExists(username string, excludeID int) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE username = $1 AND id != $2`

	err := ur.db.Get(&count, query, username, excludeID)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// UserService untuk business logic
type UserService struct {
	repo   *UserRepository
	logger *logrus.Logger
}

// NewUserService membuat instance baru UserService
func NewUserService(repo *UserRepository, logger *logrus.Logger) *UserService {
	return &UserService{
		repo:   repo,
		logger: logger,
	}
}

// CreateUser membuat user baru dengan validasi
func (us *UserService) CreateUser(req CreateUserRequest) (*User, error) {
	// Check if email already exists
	emailExists, err := us.repo.CheckEmailExists(req.Email, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if emailExists {
		return nil, fmt.Errorf("email already exists")
	}

	// Check if username already exists
	usernameExists, err := us.repo.CheckUsernameExists(req.Username, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to check username: %w", err)
	}
	if usernameExists {
		return nil, fmt.Errorf("username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		us.logger.WithError(err).Error("Failed to hash password")
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	now := time.Now()
	user := &User{
		Username:  req.Username,
		Email:     req.Email,
		FullName:  req.FullName,
		Password:  string(hashedPassword),
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := us.repo.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID mengambil user berdasarkan ID
func (us *UserService) GetUserByID(id int) (*User, error) {
	return us.repo.GetUserByID(id)
}

// GetAllUsers mengambil semua users dengan pagination
func (us *UserService) GetAllUsers(page, limit int) ([]User, int, error) {
	offset := (page - 1) * limit
	return us.repo.GetAllUsers(limit, offset)
}

// UpdateUser mengupdate user dengan validasi
func (us *UserService) UpdateUser(id int, req UpdateUserRequest) error {
	updates := make(map[string]interface{})

	if req.Username != "" {
		usernameExists, err := us.repo.CheckUsernameExists(req.Username, id)
		if err != nil {
			return fmt.Errorf("failed to check username: %w", err)
		}
		if usernameExists {
			return fmt.Errorf("username already exists")
		}
		updates["username"] = req.Username
	}

	if req.Email != "" {
		emailExists, err := us.repo.CheckEmailExists(req.Email, id)
		if err != nil {
			return fmt.Errorf("failed to check email: %w", err)
		}
		if emailExists {
			return fmt.Errorf("email already exists")
		}
		updates["email"] = req.Email
	}

	if req.FullName != "" {
		updates["full_name"] = req.FullName
	}

	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	return us.repo.UpdateUser(id, updates)
}

// DeleteUser menghapus user
func (us *UserService) DeleteUser(id int) error {
	return us.repo.DeleteUser(id)
}

// UserHandler untuk HTTP handlers
type UserHandler struct {
	service *UserService
	logger  *logrus.Logger
}

// NewUserHandler membuat instance baru UserHandler
func NewUserHandler(service *UserService, logger *logrus.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

// CreateUser handler untuk POST /users
func (uh *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		uh.logger.WithError(err).Error("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid JSON format",
			"message": err.Error(),
		})
		return
	}

	user, err := uh.service.CreateUser(req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "email already exists" || err.Error() == "username already exists" {
			statusCode = http.StatusConflict
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"error":   "Failed to create user",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "User created successfully",
		"data":    user,
	})
}

// GetUsers handler untuk GET /users
func (uh *UserHandler) GetUsers(c *gin.Context) {
	page := 1
	limit := 10

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	users, total, err := uh.service.GetAllUsers(page, limit)
	if err != nil {
		uh.logger.WithError(err).Error("Failed to get users")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get users",
			"message": err.Error(),
		})
		return
	}

	totalPages := (total + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Users retrieved successfully",
		"data":    users,
		"meta": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// GetUser handler untuk GET /users/:id
func (uh *UserHandler) GetUser(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid user ID format",
			"message": "User ID must be a number",
		})
		return
	}

	user, err := uh.service.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "User not found",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User retrieved successfully",
		"data":    user,
	})
}

// UpdateUser handler untuk PUT /users/:id
func (uh *UserHandler) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid user ID format",
			"message": "User ID must be a number",
		})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		uh.logger.WithError(err).Error("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid JSON format",
			"message": err.Error(),
		})
		return
	}

	err = uh.service.UpdateUser(id, req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "email already exists" || err.Error() == "username already exists" {
			statusCode = http.StatusConflict
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"error":   "Failed to update user",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User updated successfully",
	})
}

// DeleteUser handler untuk DELETE /users/:id
func (uh *UserHandler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid user ID format",
			"message": "User ID must be a number",
		})
		return
	}

	err = uh.service.DeleteUser(id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"error":   "Failed to delete user",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User deleted successfully",
	})
}

// HealthCheck handler untuk GET /health
func (uh *UserHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   "user-management-service",
		"timestamp": time.Now().UTC(),
		"version":   "1.0.0",
	})
}

// initDatabase inisialisasi database dan create table
func initDatabase() (*sqlx.DB, error) {
	// Database connection string
	dbURL := "postgres://postgres:password@localhost:5432/microservices_db?sslmode=disable"

	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Create users table jika belum ada
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			full_name VARCHAR(100) NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
		CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
		CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);
	`

	if _, err := db.Exec(createTableQuery); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return db, nil
}

func main() {
	// Setup logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	logger.Info("Starting User Management Service...")

	// Initialize database
	db, err := initDatabase()
	if err != nil {
		logger.WithError(err).Fatal("Failed to initialize database")
	}
	defer db.Close()

	logger.Info("Database connected successfully")

	// Setup repository, service, dan handler
	userRepo := NewUserRepository(db, logger)
	userService := NewUserService(userRepo, logger)
	userHandler := NewUserHandler(userService, logger)

	// Setup Gin router
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Routes
	api := router.Group("/api/v1")
	{
		api.GET("/health", userHandler.HealthCheck)
		api.POST("/users", userHandler.CreateUser)
		api.GET("/users", userHandler.GetUsers)
		api.GET("/users/:id", userHandler.GetUser)
		api.PUT("/users/:id", userHandler.UpdateUser)
		api.DELETE("/users/:id", userHandler.DeleteUser)
	}

	// Setup server
	server := &http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	// Start server dalam goroutine
	go func() {
		logger.WithField("port", 8081).Info("Server starting...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("Failed to start server")
		}
	}()

	// Wait for interrupt signal untuk graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown dengan timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.WithError(err).Error("Server forced to shutdown")
	} else {
		logger.Info("Server shutdown completed")
	}
}

// Helper function untuk join strings
func joinStrings(strings []string, separator string) string {
	if len(strings) == 0 {
		return ""
	}

	result := strings[0]
	for i := 1; i < len(strings); i++ {
		result += separator + strings[i]
	}
	return result
}
