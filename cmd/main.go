package main

import (
	"log"
	"net/http"
	"os"

	"github.com/RezaBG/Inventory-management-api/internal/middleware"
	"github.com/RezaBG/Inventory-management-api/internal/platform/db"
	"github.com/RezaBG/Inventory-management-api/internal/product"
	"github.com/RezaBG/Inventory-management-api/internal/supplier"
	"github.com/RezaBG/Inventory-management-api/internal/user"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: .env file not found, using default environment variables")
	}

	// Connect to the database
	database, err := db.ConnectDatabase()
	if err != nil {
		log.Fatalf("Fatal error: could not connect to database: %v", err)
	}

	// --- Run database migrations ---
	log.Println("Running database migrations...")
	err = database.AutoMigrate(
		&product.Product{},
		&user.User{},
		&user.RefreshToken{},
		&supplier.Supplier{},
	)
	if err != nil {
		log.Fatalf("Fatal error: could not run migrations: %v", err)
	}
	log.Println("Database migrations completed successfully.")

	// --- Dependency Injection ---
	// User Feature
	userRepo := user.NewRepository(database)
	refreshTokenRepo := user.NewRefreshTokenRepository(database)
	userSvc := user.NewService(userRepo, refreshTokenRepo)
	userHandler := user.NewHandler(userSvc)

	// Product Feature
	productRepo := product.NewRepository(database)
	productSvc := product.NewService(productRepo)
	productHandler := product.NewHandler(productSvc)

	// Supplier Feature
	supplierRepo := supplier.NewRepository(database)
	supplierSvc := supplier.NewService(supplierRepo)
	supplierHandler := supplier.NewHandler(supplierSvc)

	// --- Middleware ---
	authMiddleware := middleware.AuthMiddleware(userSvc)

	// --- Setup Router ---
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default port
	}
	router := gin.Default()

	// --- Register Routes ---
	router.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})
	})
	// Register the public /login and /register routes.
	// product.RegisterRoutes(router, productHandler)
	user.RegisterAuthRoutes(router, userHandler)

	// Protected Routes (Requires a valid JWT)
	protectedRoutes := router.Group("/")
	protectedRoutes.Use(authMiddleware)
	{
		product.RegisterRoutes(protectedRoutes, productHandler)
		supplier.RegisterRoutes(protectedRoutes, supplierHandler)
	}

	// --- Start Server ---
	log.Printf("Server is running on port %s", port)
	router.Run(":" + port) // listen and serve on ":8080"
}
