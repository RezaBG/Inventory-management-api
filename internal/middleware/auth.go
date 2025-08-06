package middleware

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/RezaBG/Inventory-management-api/internal/user"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(userSvc user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHandler := c.GetHeader("Authorization")
		if authHandler == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "Authorization header is required"},
			)
			return
		}

		// The header should be in the format "Bearer <token>"
		parts := strings.Split(authHandler, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "Invalid authorization format"},
			)
			return
		}

		tokenString := parts[1]
		claims := &jwt.RegisteredClaims{}

		token, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			func(token *jwt.Token) (interface{}, error) {
				// We use HMAC, so we need to check the signing method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}
				return []byte(os.Getenv("JWT_SECRET")), nil
			})

		// This is where we habdle the token validation errors you asked about
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.AbortWithStatusJSON(
					http.StatusUnauthorized,
					gin.H{"error": "Token has expired"},
				)
			} else {
				c.AbortWithStatusJSON(
					http.StatusUnauthorized,
					gin.H{"error": "Invalid token"},
				)
			}
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "Token is not valid"},
			)
			return
		}

		// Token is valid, let's get the user ID from "subject" claim
		userID, err := strconv.Atoi(claims.Subject)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "Invalid user ID in token"},
			)
			return
		}

		foundUser, err := userSvc.FindByID(uint(userID))

		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "User not found"},
			)
			return
		}

		// Set the full user object in the context
		c.Set("currentUser", foundUser)

		// Call the next handler in the chain
		c.Next()
	}
}
