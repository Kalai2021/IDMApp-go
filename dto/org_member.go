package dto

import (
	"github.com/google/uuid"
	"time"
)

type OrgMemberOpRequest struct {
	Op       int       `json:"op" binding:"required"` // 1 for ADD, 2 for REMOVE
	Type     string    `json:"type" binding:"required"` // USER, GROUP, or ROLE
	OrgID    uuid.UUID `json:"orgId" binding:"required"`
	EntityID uuid.UUID `json:"entityId" binding:"required"`
}

type OrgMemberResponse struct {
	OrgID     uuid.UUID `json:"orgId"`
	EntityID  uuid.UUID `json:"entityId"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"createdAt"`
} 