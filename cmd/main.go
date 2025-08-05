package main

import (
	"log"
	"net/http"
	"os"

	"github.com/RezaBG/Inventory-management-api/internal/platform/db"
	"github.com/RezaBG/Inventory-management-api/internal/product"

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

	// Run database migrations
	log.Println("Running database migrations...")
	err = database.AutoMigrate(&product.Product{})
	if err != nil {
		log.Fatalf("Fatal error: could not run migrations: %v", err)
	}
	log.Println("Database migrations completed successfully.")

	// Dependency Injection wiring
	productRepository := product.NewRepository(database)
	productService := product.NewService(productRepository)
	productHandler := product.NewHandler(productService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default port
	}

	router := gin.Default()

	product.RegisterRoutes(router, productHandler)

	router.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})
	})

	log.Printf("Server is running on port %s", port)
	router.Run(":" + port) // listen and serve on ":8080"
}
