package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserGroupMember struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	GroupID   uuid.UUID `json:"groupId" gorm:"type:uuid;not null;column:group_id"`
	UserID    uuid.UUID `json:"userId" gorm:"type:uuid;not null;column:user_id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (m *UserGroupMember) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

func (m *UserGroupMember) TableName() string {
	return "user_group_members"
} 