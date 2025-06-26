package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrgMember struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrgID     uuid.UUID `json:"orgId" gorm:"type:uuid;not null;column:org_id"`
	EntityID  uuid.UUID `json:"entityId" gorm:"type:uuid;not null;column:entity_id"`
	Type      string    `json:"type" gorm:"type:varchar(32);not null"` // USER, GROUP, or ROLE
	CreatedAt time.Time `json:"createdAt"`
}

func (m *OrgMember) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

func (m *OrgMember) TableName() string {
	return "org_members"
} 