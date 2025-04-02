package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	UserID      uint     `gorm:"primaryKey"`
	Username    string   `gorm:"uniqueIndex;not null"`
	Email       string   `gorm:"uniqueIndex;not null"`
	Password    string   `gorm:"not null"`                // Stored as bcrypt hash
	Role        string   `gorm:"not null;default:'user'"` // 'admin' or 'user'
	TOTPSecret  string   `gorm:"type:varchar(32)"`        // TOTP secret key
	TOTPEnabled bool     `gorm:"default:false"`           // Whether 2FA is enabled
	BackupCodes []string `gorm:"type:text[]"`             // Backup codes for 2FA
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// UserLoginRequest represents the login request
type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	TOTPCode string `json:"totp_code,omitempty"`
}

// UserRegisterRequest represents the registration request
type UserRegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// UserResponse represents the user data returned to the client
type UserResponse struct {
	UserID      uint      `json:"user_id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Role        string    `json:"role"`
	TOTPEnabled bool      `json:"totp_enabled"`
	CreatedAt   time.Time `json:"created_at"`
}

// TOTPSetupResponse represents the response when setting up 2FA
type TOTPSetupResponse struct {
	Secret      string   `json:"secret"`
	QRCode      string   `json:"qr_code"`
	BackupCodes []string `json:"backup_codes"`
}

// ToResponse converts a User to a UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		UserID:      u.UserID,
		Username:    u.Username,
		Email:       u.Email,
		Role:        u.Role,
		TOTPEnabled: u.TOTPEnabled,
		CreatedAt:   u.CreatedAt,
	}
}
