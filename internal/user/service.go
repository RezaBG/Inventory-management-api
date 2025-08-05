package user

import (
	"errors"
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (s *service) CreateNewUser(input CreateUserInput) (*User, error) {
	_, err := s.repo.FindByEmail(input.Email)
	if err == nil {
		return nil, fmt.Errorf("email already exists")
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
