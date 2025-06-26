package role

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RoleService struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewRoleService(db *gorm.DB) *RoleService {
	return &RoleService{
		db:     db,
		logger: logrus.New(),
	}
}

func (s *RoleService) GetAllRoles() ([]Role, error) {
	var roles []Role
	result := s.db.Find(&roles)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get roles: %w", result.Error)
	}
	return roles, nil
}

func (s *RoleService) GetRole(id uuid.UUID) (*Role, error) {
	var role Role
	result := s.db.First(&role, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get role: %w", result.Error)
	}
	return &role, nil
}

func (s *RoleService) CreateRole(req RoleCreateRequest) (*Role, error) {
	// Check if role with name already exists
	var existingRole Role
	if err := s.db.Where("name = ?", req.Name).First(&existingRole).Error; err == nil {
		return nil, errors.New("role with this name already exists")
	}

	role := Role{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	result := s.db.Create(&role)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create role: %w", result.Error)
	}

	return &role, nil
}

func (s *RoleService) UpdateRole(id uuid.UUID, req RoleUpdateRequest) (*Role, error) {
	var role Role
	result := s.db.First(&role, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get role: %w", result.Error)
	}

	// Update fields if provided
	if req.Name != "" {
		// Check if name is already taken by another role
		var existingRole Role
		if err := s.db.Where("name = ? AND id != ?", req.Name, id).First(&existingRole).Error; err == nil {
			return nil, errors.New("role name is already taken by another role")
		}
		role.Name = req.Name
	}
	if req.DisplayName != "" {
		role.DisplayName = req.DisplayName
	}
	if req.Description != "" {
		role.Description = req.Description
	}

	role.UpdatedAt = time.Now()

	result = s.db.Save(&role)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to update role: %w", result.Error)
	}

	return &role, nil
}

func (s *RoleService) DeleteRole(id uuid.UUID) error {
	result := s.db.Delete(&Role{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete role: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("role not found")
	}
	return nil
}
