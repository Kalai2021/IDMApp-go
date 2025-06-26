package dto

import (
	"github.com/google/uuid"
	"time"
)

type RoleMemberOpRequest struct {
	Op       int       `json:"op" binding:"required"` // 1 for ADD, 2 for REMOVE
	Type     string    `json:"type" binding:"required"` // USER or GROUP
	RoleID   uuid.UUID `json:"roleId" binding:"required"`
	EntityID uuid.UUID `json:"entityId" binding:"required"`
}

type RoleMemberResponse struct {
	RoleID    uuid.UUID `json:"roleId"`
	EntityID  uuid.UUID `json:"entityId"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"createdAt"`
} 