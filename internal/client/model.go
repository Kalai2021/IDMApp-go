package client

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Client struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ClientID     string         `gorm:"uniqueIndex;not null"`
	ClientSecret string         `gorm:"not null"`
	Name         string         `gorm:"not null"`
	RedirectURIs pq.StringArray `gorm:"type:text[]"`
	Scopes       pq.StringArray `gorm:"type:text[];not null"`
	Active       bool           `gorm:"not null;default:true"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
