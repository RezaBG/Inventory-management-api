package main

import (
	"log"
	"net/http"
	"os"

	// ADDED: Imports for Swagger documentation
	_ "github.com/RezaBG/Inventory-management-api/docs" // This links to the generated docs.
	"github.com/RezaBG/Inventory-management-api/internal/inventory"
	"github.com/RezaBG/Inventory-management-api/internal/middleware"
	"github.com/RezaBG/Inventory-management-api/internal/platform/db"
	"github.com/RezaBG/Inventory-management-api/internal/product"
	"github.com/RezaBG/Inventory-management-api/internal/supplier"
	"github.com/RezaBG/Inventory-management-api/internal/user"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Inventory Management API
// @version         1.0
// @description     This is a server for an inventory management system, built with Go and Gin.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Reza Barzegar
// @contact.url    https://github.com/RezaBG

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:2019
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and a valid JWT token.
func main() {
	// ... (godotenv, db connection, migrations, and DI are unchanged) ...
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	database, err := db.ConnectDatabase()
	if err != nil {
		log.Fatalf("Fatal error: could not connect to database: %v", err)
	}

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
		port = "8080"
	}
	router := gin.Default()

	// --- Register Routes ---

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	user.RegisterAuthRoutes(router, userHandler)
	router.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
	})

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
	router.Run(":" + port)
}
