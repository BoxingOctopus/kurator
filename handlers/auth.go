package handlers

import (
	"media_collection/auth"
	"media_collection/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var DB *gorm.DB

// SetDB sets the database connection
func SetDB(db *gorm.DB) {
	DB = db
}

// SetupTOTP handles the 2FA setup process
func SetupTOTP(c *gin.Context) {
	userID := c.GetUint("user_id")
	var user models.User
	if err := DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Generate TOTP secret
	secret, err := auth.GenerateTOTPSecret(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate TOTP secret"})
		return
	}

	// Generate QR code
	qrCode, err := auth.GenerateTOTPQRCode(user.Username, secret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR code"})
		return
	}

	// Generate backup codes
	backupCodes, err := auth.GenerateBackupCodes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate backup codes"})
		return
	}

	c.JSON(http.StatusOK, models.TOTPSetupResponse{
		Secret:      secret,
		QRCode:      qrCode,
		BackupCodes: backupCodes,
	})
}

// EnableTOTP enables 2FA for a user
func EnableTOTP(c *gin.Context) {
	userID := c.GetUint("user_id")
	var user models.User
	if err := DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Validate TOTP code
	if !auth.ValidateTOTP(user.TOTPSecret, req.Code) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid TOTP code"})
		return
	}

	// Enable 2FA
	user.TOTPEnabled = true
	if err := DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enable 2FA"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "2FA enabled successfully"})
}

// DisableTOTP disables 2FA for a user
func DisableTOTP(c *gin.Context) {
	userID := c.GetUint("user_id")
	var user models.User
	if err := DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Validate TOTP code
	if !auth.ValidateTOTP(user.TOTPSecret, req.Code) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid TOTP code"})
		return
	}

	// Disable 2FA
	user.TOTPEnabled = false
	user.TOTPSecret = ""
	user.BackupCodes = nil
	if err := DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to disable 2FA"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "2FA disabled successfully"})
}
