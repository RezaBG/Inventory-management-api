package user

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service interface {
	CreateNewUser(input CreateUserInput) (*User, error)
	Login(input LoginInput) (*LoginResponse, error)
	FindByID(id uint) (*User, error)
	RefreshToken(input RefreshTokenInput) (*AccessTokenResponse, error)
}

type service struct {
	userRepo Repository
	rtRepo   RefreshTokenRepository
}

func NewService(userRepo Repository, rtRepo RefreshTokenRepository) Service {
	return &service{
		userRepo: userRepo,
		rtRepo:   rtRepo,
	}
}

func (s *service) CreateNewUser(input CreateUserInput) (*User, error) {
	_, err := s.userRepo.FindByEmail(input.Email)
	if err == nil {
		return nil, fmt.Errorf("email already in use")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("database error: %w", err)
	}

	if err := validatePassword(input.Password); err != nil {
		return nil, err
	}

	// Step 2 will be hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	newUser := User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	if err := s.userRepo.Save(&newUser); err != nil {
		return nil, err
	}

	return &newUser, nil
}

func (s *service) Login(input LoginInput) (*LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(input.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Refresh token logic
	refreshTokenString, err := generateSecureRandomToken(32)
	if err != nil {
		return nil, fmt.Errorf("could not generate refresh token: %w", err)
	}

	rtExpirationHours, _ := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRATION_HOURS"))
	if rtExpirationHours == 0 {
		rtExpirationHours = 168 // Default to 7 days
	}

	refreshToken := &RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenString,
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(rtExpirationHours)),
	}

	if err := s.rtRepo.Create(refreshToken); err != nil {
		return nil, fmt.Errorf("could not save refresh token: %w", err)
	}

	atExpirationMinutes, _ := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXPIRATION_MINUTES"))
	if atExpirationMinutes == 0 {
		atExpirationMinutes = 24
	}

	claims := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    os.Getenv("JWT_ISSUER"),
			Subject:   strconv.FormatUint(uint64(user.ID), 10),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(atExpirationMinutes))),
		})

	// Sign the token with a secret key (loaded from environment variable)
	accessTokenString, err := claims.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, fmt.Errorf("could not create token: %w", err)
	}

	return &LoginResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil

}

func (s *service) RefreshToken(input RefreshTokenInput) (*AccessTokenResponse, error) {
	refreshToken, err := s.rtRepo.FindByToken(input.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	if time.Now().After(refreshToken.ExpiresAt) {
		return nil, fmt.Errorf("refresh token has expired")
	}

	atExpirationMinutes, _ := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXPIRATION_MINUTES"))
	if atExpirationMinutes == 0 {
		atExpirationMinutes = 15
	}

	claims := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    os.Getenv("JWT_ISSUER"),
			Subject:   strconv.FormatUint(uint64(refreshToken.UserID), 10),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(atExpirationMinutes))),
		})

	newAccessTokenString, err := claims.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, fmt.Errorf("could not create new access token: %w", err)
	}

	return &AccessTokenResponse{AccessToken: newAccessTokenString}, nil
}

func generateSecureRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// FundByID retrieves a user by their ID.
func (s *service) FindByID(id uint) (*User, error) {
	return s.userRepo.FindByID(id)
}

func validatePassword(password string) error {
	if ok, _ := regexp.MatchString(`[A-Z]`, password); !ok {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}

	if ok, _ := regexp.MatchString(`[a-z]`, password); !ok {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}

	if ok, _ := regexp.MatchString(`[0-9]`, password); !ok {
		return fmt.Errorf("password must contain at least one digit")
	}

	if ok, _ := regexp.MatchString(`[\W]`, password); !ok {
		return fmt.Errorf("password must contain at least one special character")
	}

	return nil // Password is valid
}
