package auth

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"math/big"
	"strings"

	"github.com/pquerna/otp/totp"
)

const (
	// TOTPIssuer is the issuer name for TOTP
	TOTPIssuer = "Media Collection"
	// BackupCodeLength is the length of each backup code
	BackupCodeLength = 10
	// NumBackupCodes is the number of backup codes to generate
	NumBackupCodes = 8
)

// GenerateTOTPSecret generates a new TOTP secret for a user
func GenerateTOTPSecret(username string) (string, error) {
	// Generate a random secret
	secret := make([]byte, 20)
	if _, err := rand.Read(secret); err != nil {
		return "", fmt.Errorf("failed to generate secret: %w", err)
	}

	// Encode the secret in base32
	secretBase32 := base32.StdEncoding.EncodeToString(secret)
	return secretBase32, nil
}

// GenerateTOTPQRCode generates a QR code for TOTP setup
func GenerateTOTPQRCode(username, secret string) (string, error) {
	// Decode the base32 secret
	secretBytes, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", fmt.Errorf("failed to decode secret: %w", err)
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      TOTPIssuer,
		AccountName: username,
		Secret:      secretBytes,
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate TOTP key: %w", err)
	}

	return key.URL(), nil
}

// ValidateTOTP validates a TOTP code
func ValidateTOTP(secret, code string) bool {
	return totp.Validate(code, secret)
}

// GenerateBackupCodes generates backup codes for 2FA
func GenerateBackupCodes() ([]string, error) {
	codes := make([]string, NumBackupCodes)
	for i := 0; i < NumBackupCodes; i++ {
		// Generate a random number
		num, err := rand.Int(rand.Reader, big.NewInt(1<<40))
		if err != nil {
			return nil, fmt.Errorf("failed to generate backup code: %w", err)
		}

		// Convert to base32 and pad with zeros
		code := base32.StdEncoding.EncodeToString(num.Bytes())
		code = strings.Repeat("0", BackupCodeLength-len(code)) + code
		codes[i] = code
	}
	return codes, nil
}

// ValidateBackupCode validates a backup code
func ValidateBackupCode(code string, backupCodes []string) bool {
	for i, storedCode := range backupCodes {
		if code == storedCode {
			// Remove the used backup code
			backupCodes[i] = backupCodes[len(backupCodes)-1]
			backupCodes = backupCodes[:len(backupCodes)-1]
			return true
		}
	}
	return false
}
