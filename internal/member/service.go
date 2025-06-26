package member

import (
	"errors"
	"fmt"
	"time"

	"idmapp-go/dto"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type MemberService struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewMemberService(db *gorm.DB) *MemberService {
	return &MemberService{
		db:     db,
		logger: logrus.New(),
	}
}

func (s *MemberService) GetMember(groupID, userID uuid.UUID) (*Member, error) {
	if groupID == uuid.Nil || userID == uuid.Nil {
		return nil, errors.New("group ID and user ID cannot be null")
	}

	var member Member
	result := s.db.Where("group_id = ? AND user_id = ?", groupID, userID).First(&member)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get member: %w", result.Error)
	}
	return &member, nil
}

func (s *MemberService) GetAllMembers() ([]Member, error) {
	var members []Member
	result := s.db.Find(&members)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get members: %w", result.Error)
	}
	return members, nil
}

func (s *MemberService) GetMembersByGroupID(groupID uuid.UUID) ([]Member, error) {
	if groupID == uuid.Nil {
		return nil, errors.New("group ID cannot be null")
	}

	var members []Member
	result := s.db.Where("group_id = ?", groupID).Find(&members)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get members by group ID: %w", result.Error)
	}
	return members, nil
}

func (s *MemberService) GetMembersByUserID(userID uuid.UUID) ([]Member, error) {
	if userID == uuid.Nil {
		return nil, errors.New("user ID cannot be null")
	}

	var members []Member
	result := s.db.Where("user_id = ?", userID).Find(&members)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get members by user ID: %w", result.Error)
	}
	return members, nil
}

func (s *MemberService) AddMember(req dto.MemberOpRequest) (*Member, error) {
	// Check if member already exists
	existingMember, err := s.GetMember(req.GroupID, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing member: %w", err)
	}
	if existingMember != nil {
		return nil, errors.New("member already exists in this group")
	}

	member := Member{
		GroupID:   req.GroupID,
		UserID:    req.UserID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result := s.db.Create(&member)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to add member: %w", result.Error)
	}

	return &member, nil
}

func (s *MemberService) RemoveMember(groupID, userID uuid.UUID) error {
	if groupID == uuid.Nil || userID == uuid.Nil {
		return errors.New("group ID and user ID cannot be null")
	}

	result := s.db.Where("group_id = ? AND user_id = ?", groupID, userID).Delete(&Member{})
	if result.Error != nil {
		return fmt.Errorf("failed to remove member: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("member not found")
	}
	return nil
}

func (s *MemberService) ProcessMemberOperation(req dto.MemberOpRequest) (*Member, error) {
	switch req.Op {
	case dto.OpTypeAdd:
		return s.AddMember(req)
	case dto.OpTypeRemove:
		err := s.RemoveMember(req.GroupID, req.UserID)
		if err != nil {
			return nil, err
		}
		return nil, nil // Return nil for successful removal
	default:
		return nil, fmt.Errorf("invalid operation type: %s", req.Op)
	}
}
