package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"social-network/pkg/models"
)

// validate all fields when user registers
func ValidateNewUser(user models.User) error {
	if err := validateFirstName(user.FirstName); err != nil {
		return fmt.Errorf("first name: %w", err)
	}
	if err := validateLastName(user.LastName); err != nil {
		return fmt.Errorf("last name: %w", err)
	}
	if err := validateBirth(user.DateOfBirth); err != nil {
		return fmt.Errorf("date of birth: %w", err)
	}
	if err := validatePassword(user.Password); err != nil {
		return fmt.Errorf("password: %w", err)
	}
	if err := validateEmail(user.Email); err != nil {
		return fmt.Errorf("email: %w", err)
	}
	return nil
}

func validateFirstName(name string) error {
	name = strings.TrimSpace(name)
	if fieldEmpty(name) {
		return errors.New("cannot be empty")
	}
	if len(name) > 50 {
		return errors.New("too long")
	}
	if !regexp.MustCompile(`^[a-zA-Z\-']+$`).MatchString(name) {
		return errors.New("invalid characters")
	}
	return nil
}

func validateLastName(name string) error {
	name = strings.TrimSpace(name)
	if fieldEmpty(name) {
		return errors.New("cannot be empty")
	}
	if len(name) > 50 {
		return errors.New("too long")
	}
	if !regexp.MustCompile(`^[a-zA-Z\-']+$`).MatchString(name) {
		return errors.New("invalid characters")
	}
	return nil
}

func validateBirth(birthDate string) error {
	birthDate = strings.TrimSpace(birthDate)
	if fieldEmpty(birthDate) {
		return errors.New("cannot be empty")
	}
	if err := validateDate(birthDate); err != nil {
		return fmt.Errorf("invalid date format: %w", err)
	}
	parsed, _ := time.Parse("2006-01-02", birthDate)
	if parsed.After(time.Now()) {
		return errors.New("birth date cannot be in the future")
	}
	return nil
}

func validatePassword(password string) error {
	password = strings.TrimSpace(password)
	if fieldEmpty(password) {
		return errors.New("cannot be empty")
	}
	if len(password) < 8 {
		return errors.New("must be at least 8 characters")
	}
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return errors.New("must contain an uppercase letter")
	}
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return errors.New("must contain a lowercase letter")
	}
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return errors.New("must contain a number")
	}
	if !regexp.MustCompile(`[!@#\$%\^&\*]`).MatchString(password) {
		return errors.New("must contain a special character (!@#$%^&*)")
	}
	return nil
}

func validateEmail(email string) error {
	email = strings.TrimSpace(email)
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if fieldEmpty(email) {
		return errors.New("cannot be empty")
	}
	if len(email) > 100 {
		return errors.New("too long")
	}
	if !emailRegex.MatchString(email) {
		return errors.New("invalid format")
	}
	return nil
}

func validateDate(dateStr string) error {
	layout := "2006-01-02"
	_, err := time.Parse(layout, dateStr)
	if err != nil {
		return fmt.Errorf("should be YYYY-MM-DD: %s", err)
	}
	return nil
}

func fieldEmpty(value string) bool {
	return len(value) == 0
}
