package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// CreateUser handles the API request to create a new user.
// @Summary      Register a new user
// @Description  Creates a new user account with the provided details.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        user body CreateUserInput true "User Registration Info"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Router       /register [post]
func (h *Handler) CreateUser(c *gin.Context) {
	var input CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdUser, err := h.svc.CreateNewUser(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"userID":  createdUser.ID,
	})
}

// Login handles the API request for user login.
// @Summary      Log in a user
// @Description  Authenticates a user and returns an access token and refresh token.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        credentials body LoginInput true "User Login Credentials"
// @Success      200  {object}  LoginResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Router       /login [post]
func (h *Handler) Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loginResponse, err := h.svc.Login(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// On success, we return the token
	c.JSON(http.StatusOK, gin.H{"token": loginResponse})
}

// RefreshToken handles the API request to get a new access token.
// @Summary      Refresh access token
// @Description  Issues a new access token in exchange for a valid refresh token.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        token body RefreshTokenInput true "Refresh Token"
// @Success      200  {object}  AccessTokenResponse
// @Failure      400  {object}  map[string]interface{} "{"error": "Error message"}"
// @Failure      401  {object}  map[string]interface{} "{"error": "Invalid refresh token"}"
// @Router       /refresh_token [post]
func (h *Handler) RefreshToken(c *gin.Context) {
	var input RefreshTokenInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.svc.RefreshToken(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
