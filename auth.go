package main

import (
	"errors"
	"time"

	"media_collection/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	UserID    uint   `gorm:"primaryKey"`
	Username  string `gorm:"uniqueIndex;not null"`
	Email     string `gorm:"uniqueIndex;not null"`
	Password  string `gorm:"not null"`                // Stored as bcrypt hash
	Role      string `gorm:"not null;default:'user'"` // 'admin' or 'user'
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserLoginRequest represents the login request
type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserRegisterRequest represents the registration request
type UserRegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// UserResponse represents the user data returned to the client
type UserResponse struct {
	UserID    uint      `json:"user_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// HashPassword securely hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPassword checks if the provided password matches the hash
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// CreateUser creates a new user in the database
func CreateUser(user *models.User) error {
	// Check if username or email already exists
	var count int64
	db.Model(&models.User{}).Where("username = ? OR email = ?", user.Username, user.Email).Count(&count)
	if count > 0 {
		return errors.New("username or email already exists")
	}

	// Hash the password
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	// Set default role if not specified
	if user.Role == "" {
		user.Role = "user"
	}

	// Create the user
	return db.Create(user).Error
}

// GetUserByUsername retrieves a user by username
func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByID retrieves a user by ID
func GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates a user's information
func UpdateUser(user *models.User) error {
	// If password is being updated, hash it
	if user.Password != "" {
		hashedPassword, err := HashPassword(user.Password)
		if err != nil {
			return err
		}
		user.Password = hashedPassword
	}

	return db.Save(user).Error
}

// DeleteUser deletes a user by ID
func DeleteUser(id uint) error {
	return db.Delete(&models.User{}, id).Error
}

// ListUsers retrieves all users with pagination
func ListUsers(limit, offset int) ([]models.User, error) {
	var users []models.User
	err := db.Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}

// ToUserResponse converts a User to a UserResponse
func (u *User) ToUserResponse() UserResponse {
	return UserResponse{
		UserID:    u.UserID,
		Username:  u.Username,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
	}
}
