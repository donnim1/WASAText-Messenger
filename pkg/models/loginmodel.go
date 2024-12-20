package models

import (
	"errors"
	"fmt"
	"github.com/WASAText/db"

	"gorm.io/gorm"
)

// User represents the users table in the database
type User struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"unique;size:16" json:"name"`
}

// GetOrCreateUser retrieves a user by name or creates a new one if not found
func GetOrCreateUser(name string) (string, error) {
	var user User
	result := db.DB.Where("name = ?", name).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Create a new user
		user = User{Name: name}
		if err := db.DB.Create(&user).Error; err != nil {
			return "", fmt.Errorf("failed to create user: %w", err)
		}
	} else if result.Error != nil {
		return "", fmt.Errorf("failed to query user: %w", result.Error)
	}

	return fmt.Sprintf("%012d", user.ID), nil
}
