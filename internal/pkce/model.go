package pkce

import (
	"time"

	"github.com/google/uuid"
)

type PKCECode struct {
	ID                  uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Code                string    `gorm:"uniqueIndex;not null"`
	CodeChallenge       string    `gorm:"not null"`
	CodeChallengeMethod string    `gorm:"not null"`
	CodeVerifier        string    `gorm:"not null"`
	ClientID            string    `gorm:"not null"`
	RedirectURI         string    `gorm:"not null"`
	State               *string   `gorm:"default:null"`
	UserID              *uuid.UUID
	ExpiresAt           time.Time `gorm:"not null"`
	Used                bool      `gorm:"not null;default:false"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
