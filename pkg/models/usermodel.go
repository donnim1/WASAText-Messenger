package models

import (
	"fmt"
	"WASA/WASAText/pkg/db"
)

// User represents the users table
type User struct {
	UserID   string `gorm:"primaryKey;size:36" json:"userId"`
	Name     string `gorm:"unique;size:16" json:"name"`
	PhotoURL string `gorm:"size:255" json:"photoUrl"`
	Status   string `gorm:"size:7" json:"status"` // online, offline, away
}

// GetAllUsers retrieves all users, with optional search and pagination
func GetAllUsers(search string, limit int) ([]User, error) {
	var users []User
	query := db.DB.Limit(limit)

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve users: %w", err)
	}
	return users, nil
}

// UpdateUserName updates the user's name
func UpdateUserName(userID string, newName string) error {
	var user User
	if err := db.DB.Where("user_id = ?", userID).First(&user).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	user.Name = newName
	if err := db.DB.Save(&user).Error; err != nil {
		return fmt.Errorf("failed to update username: %w", err)
	}
	return nil
}

// UpdateUserPhoto updates the user's profile photo URL
func UpdateUserPhoto(userID string, photoURL string) error {
	var user User
	if err := db.DB.Where("user_id = ?", userID).First(&user).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	user.PhotoURL = photoURL
	if err := db.DB.Save(&user).Error; err != nil {
		return fmt.Errorf("failed to update photo: %w", err)
	}
	return nil
}
