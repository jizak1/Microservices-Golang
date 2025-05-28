package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Product model sederhana untuk contoh
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
	CreatedAt   string  `json:"created_at"`
}

// ProductService untuk business logic
type ProductService struct {
	products []Product
	logger   *logrus.Logger
}

// NewProductService membuat instance baru ProductService
func NewProductService(logger *logrus.Logger) *ProductService {
	// Data dummy untuk demo
	dummyProducts := []Product{
		{
			ID:          1,
			Name:        "Laptop Gaming",
			Description: "Laptop gaming dengan spesifikasi tinggi",
			Price:       15000000,
			Category:    "Electronics",
			CreatedAt:   time.Now().Format(time.RFC3339),
		},
		{
			ID:          2,
			Name:        "Smartphone Android",
			Description: "Smartphone Android terbaru dengan kamera canggih",
			Price:       8000000,
			Category:    "Electronics",
			CreatedAt:   time.Now().Format(time.RFC3339),
		},
		{
			ID:          3,
			Name:        "Sepatu Olahraga",
			Description: "Sepatu olahraga nyaman untuk aktivitas sehari-hari",
			Price:       750000,
			Category:    "Fashion",
			CreatedAt:   time.Now().Format(time.RFC3339),
		},
	}

	return &ProductService{
		products: dummyProducts,
		logger:   logger,
	}
}

// GetAllProducts mengambil semua produk
func (ps *ProductService) GetAllProducts() []Product {
	ps.logger.Info("Fetching all products")
	return ps.products
}

// GetProductByID mengambil produk berdasarkan ID
func (ps *ProductService) GetProductByID(id int) (*Product, error) {
	ps.logger.WithField("product_id", id).Info("Fetching product by ID")
	
	for _, product := range ps.products {
		if product.ID == id {
			return &product, nil
		}
	}
	
	return nil, fmt.Errorf("product with ID %d not found", id)
}

// AddProduct menambahkan produk baru
func (ps *ProductService) AddProduct(product Product) Product {
	// Generate ID baru (simple increment)
	maxID := 0
	for _, p := range ps.products {
		if p.ID > maxID {
			maxID = p.ID
		}
	}
	
	product.ID = maxID + 1
	product.CreatedAt = time.Now().Format(time.RFC3339)
	ps.products = append(ps.products, product)
	
	ps.logger.WithField("product_id", product.ID).Info("Product added successfully")
	return product
}

// ProductHandler untuk HTTP handlers
type ProductHandler struct {
	service *ProductService
	logger  *logrus.Logger
}

// NewProductHandler membuat instance baru ProductHandler
func NewProductHandler(service *ProductService, logger *logrus.Logger) *ProductHandler {
	return &ProductHandler{
		service: service,
		logger:  logger,
	}
}

// GetProducts handler untuk GET /products
func (ph *ProductHandler) GetProducts(c *gin.Context) {
	products := ph.service.GetAllProducts()
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Products retrieved successfully",
		"data":    products,
		"count":   len(products),
	})
}

// GetProduct handler untuk GET /products/:id
func (ph *ProductHandler) GetProduct(c *gin.Context) {
	idParam := c.Param("id")
	
	var id int
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid product ID format",
			"message": "Product ID must be a number",
		})
		return
	}
	
	product, err := ph.service.GetProductByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Product not found",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Product retrieved successfully",
		"data":    product,
	})
}

// CreateProduct handler untuk POST /products
func (ph *ProductHandler) CreateProduct(c *gin.Context) {
	var newProduct Product
	
	if err := c.ShouldBindJSON(&newProduct); err != nil {
		ph.logger.WithError(err).Error("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid JSON format",
			"message": err.Error(),
		})
		return
	}
	
	// Validasi input sederhana
	if newProduct.Name == "" || newProduct.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Validation failed",
			"message": "Name is required and price must be greater than 0",
		})
		return
	}
	
	createdProduct := ph.service.AddProduct(newProduct)
	
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Product created successfully",
		"data":    createdProduct,
	})
}

// HealthCheck handler untuk GET /health
func (ph *ProductHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   "basic-http-service",
		"timestamp": time.Now().UTC(),
		"version":   "1.0.0",
	})
}

func main() {
	// Setup logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)
	
	logger.Info("Starting Basic HTTP Service...")
	
	// Setup service dan handler
	productService := NewProductService(logger)
	productHandler := NewProductHandler(productService, logger)
	
	// Setup Gin router
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	
	// Routes
	api := router.Group("/api/v1")
	{
		api.GET("/health", productHandler.HealthCheck)
		api.GET("/products", productHandler.GetProducts)
		api.GET("/products/:id", productHandler.GetProduct)
		api.POST("/products", productHandler.CreateProduct)
	}
	
	// Setup server
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	
	// Start server dalam goroutine
	go func() {
		logger.WithField("port", 8080).Info("Server starting...")
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
