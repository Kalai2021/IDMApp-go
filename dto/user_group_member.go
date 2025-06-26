package dto

import (
	"github.com/google/uuid"
	"time"
)

type UserGroupMemberOpRequest struct {
	Op      int       `json:"op" binding:"required"` // 1 for ADD, 2 for REMOVE
	GroupID uuid.UUID `json:"groupId" binding:"required"`
	UserID  uuid.UUID `json:"userId" binding:"required"`
}

type UserGroupMemberResponse struct {
	ID        uuid.UUID `json:"id"`
	GroupID   uuid.UUID `json:"groupId"`
	UserID    uuid.UUID `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
} 