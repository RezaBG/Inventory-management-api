package main

import (
	"log"
	"net/http"
	"os"

	"github.com/RezaBG/Inventory-management-api/internal/inventory"
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
		&inventory.InventoryTransaction{},
	)
	if err != nil {
		log.Fatalf("Fatal error: could not run migrations: %v", err)
	}
	log.Println("Database migrations completed successfully.")

	// --- Dependency Injection ---
	// 1. Initialize all Repositories
	userRepo := user.NewRepository(database)
	refreshTokenRepo := user.NewRefreshTokenRepository(database)
	productRepo := product.NewRepository(database)
	supplierRepo := supplier.NewRepository(database)
	inventoryRepo := inventory.NewRepository(database)

	// 2. Initialize all Services
	userSvc := user.NewService(userRepo, refreshTokenRepo)
	supplierSvc := supplier.NewService(supplierRepo)
	productSvc := product.NewService(productRepo, inventoryRepo)
	inventorySvc := inventory.NewService(inventoryRepo, productRepo)

	// 3. Initialize all Handlers
	userHandler := user.NewHandler(userSvc)
	productHandler := product.NewHandler(productSvc)
	supplierHandler := supplier.NewHandler(supplierSvc)
	inventoryHandler := inventory.NewHandler(inventorySvc)

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
		inventory.RegisterRoutes(protectedRoutes, inventoryHandler)

	}

	// --- Start Server ---
	log.Printf("Server is running on port %s", port)
	router.Run(":" + port) // listen and serve on ":8080"
}
