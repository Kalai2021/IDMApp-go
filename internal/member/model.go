package member

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Member struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	GroupID   uuid.UUID `json:"groupId" gorm:"type:uuid;not null;column:group_id"`
	UserID    uuid.UUID `json:"userId" gorm:"type:uuid;not null;column:user_id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (m *Member) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

func (m *Member) TableName() string {
	return "members"
}
