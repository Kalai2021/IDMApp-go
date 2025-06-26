package dto

import (
	"github.com/google/uuid"
)

type OpType string

const (
	OpTypeAdd    OpType = "ADD"
	OpTypeRemove OpType = "REMOVE"
)

func (o OpType) String() string {
	return string(o)
}

func OpTypeFromString(s string) (OpType, error) {
	switch s {
	case "ADD":
		return OpTypeAdd, nil
	case "REMOVE":
		return OpTypeRemove, nil
	default:
		return "", ErrInvalidOpType
	}
}

type MemberOpRequest struct {
	Op      OpType    `json:"op" binding:"required"`
	GroupID uuid.UUID `json:"groupId" binding:"required"`
	UserID  uuid.UUID `json:"userId" binding:"required"`
}

type MemberOpResponse struct {
	ID        string `json:"id"`
	GroupID   string `json:"groupId"`
	UserID    string `json:"userId"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
