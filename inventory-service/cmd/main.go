package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Assignment/inventory-service/internal/config"
	"github.com/Assignment/inventory-service/internal/handler"
	"github.com/Assignment/inventory-service/internal/repository"
	"github.com/Assignment/inventory-service/internal/service"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.LoadConfig()

	// Connect to database
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Name,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Check DB connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Initialize repositories
	productRepo := repository.NewPostgresProductRepository(db)
	categoryRepo := repository.NewPostgresCategoryRepository(db)

	// Initialize services
	productService := service.NewProductService(productRepo)
	categoryService := service.NewCategoryService(categoryRepo)

	// Initialize router
	router := gin.Default()

	// Register routes
	api := router.Group("/api/v1")
	handler.RegisterProductRoutes(api, productService)
	handler.RegisterCategoryRoutes(api, categoryService)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Inventory Service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
