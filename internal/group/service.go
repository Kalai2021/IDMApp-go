package group

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type GroupService struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewGroupService(db *gorm.DB) *GroupService {
	return &GroupService{
		db:     db,
		logger: logrus.New(),
	}
}

func (s *GroupService) GetAllGroups() ([]Group, error) {
	var groups []Group
	result := s.db.Find(&groups)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get groups: %w", result.Error)
	}
	return groups, nil
}

func (s *GroupService) GetGroup(id uuid.UUID) (*Group, error) {
	var group Group
	result := s.db.First(&group, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get group: %w", result.Error)
	}
	return &group, nil
}

func (s *GroupService) CreateGroup(req GroupCreateRequest) (*Group, error) {
	// Check if group with name already exists
	var existingGroup Group
	if err := s.db.Where("name = ?", req.Name).First(&existingGroup).Error; err == nil {
		return nil, errors.New("group with this name already exists")
	}

	group := Group{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	result := s.db.Create(&group)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create group: %w", result.Error)
	}

	return &group, nil
}

func (s *GroupService) UpdateGroup(id uuid.UUID, req GroupUpdateRequest) (*Group, error) {
	var group Group
	result := s.db.First(&group, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get group: %w", result.Error)
	}

	// Update fields if provided
	if req.Name != "" {
		// Check if name is already taken by another group
		var existingGroup Group
		if err := s.db.Where("name = ? AND id != ?", req.Name, id).First(&existingGroup).Error; err == nil {
			return nil, errors.New("group name is already taken by another group")
		}
		group.Name = req.Name
	}
	if req.DisplayName != "" {
		group.DisplayName = req.DisplayName
	}
	if req.Description != "" {
		group.Description = req.Description
	}

	group.UpdatedAt = time.Now()

	result = s.db.Save(&group)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to update group: %w", result.Error)
	}

	return &group, nil
}

func (s *GroupService) DeleteGroup(id uuid.UUID) error {
	result := s.db.Delete(&Group{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete group: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("group not found")
	}
	return nil
}
