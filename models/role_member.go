package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleMember struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	RoleID    uuid.UUID `json:"roleId" gorm:"type:uuid;not null;column:role_id"`
	EntityID  uuid.UUID `json:"entityId" gorm:"type:uuid;not null;column:entity_id"`
	Type      string    `json:"type" gorm:"type:varchar(32);not null"` // USER or GROUP
	CreatedAt time.Time `json:"createdAt"`
}

func (m *RoleMember) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

func (m *RoleMember) TableName() string {
	return "role_members"
} 