package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/RezaBG/Inventory-management-api/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println(" DEBUG: .env values loaded successfully.")
	fmt.Println("DB_USER:", os.Getenv("DB_USER"))
	fmt.Println("DB_NAME:", os.Getenv("DB_NAME"))

	db.ConnectDatabase()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default port
	}

	r := gin.Default()

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})
	})

	r.Run(":" + port) // listen and serve on ":8080"
}
