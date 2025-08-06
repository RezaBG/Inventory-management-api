package user

import (
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
	Login(input LoginInput) (string, error)
	FindByID(id uint) (*User, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateNewUser(input CreateUserInput) (*User, error) {
	_, err := s.repo.FindByEmail(input.Email)
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

	if err := s.repo.Save(&newUser); err != nil {
		return nil, err
	}

	return &newUser, nil
}

func (s *service) Login(input LoginInput) (string, error) {
	user, err := s.repo.FindByEmail(input.Email)
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	expirationHours, _ := strconv.Atoi(os.Getenv("JWT_EXPIRATION_HOURS"))
	if expirationHours == 0 {
		expirationHours = 24
	}

	claims := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    os.Getenv("JWT_ISSUER"),
			Subject:   strconv.FormatUint(uint64(user.ID), 10),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(expirationHours))),
		})

	// Sign the token with a secret key (loaded from environment variable)
	tokenString, err := claims.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", fmt.Errorf("could not create token: %w", err)
	}

	return tokenString, nil

}

// FundByID retrieves a user by their ID.
func (s *service) FindByID(id uint) (*User, error) {
	return s.repo.FindByID(id)
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
