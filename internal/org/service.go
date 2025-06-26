package org

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OrgService struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewOrgService(db *gorm.DB) *OrgService {
	return &OrgService{
		db:     db,
		logger: logrus.New(),
	}
}

func (s *OrgService) GetAllOrgs() ([]Org, error) {
	var orgs []Org
	result := s.db.Find(&orgs)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get organizations: %w", result.Error)
	}
	return orgs, nil
}

func (s *OrgService) GetOrg(id uuid.UUID) (*Org, error) {
	var org Org
	result := s.db.First(&org, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get organization: %w", result.Error)
	}
	return &org, nil
}

func (s *OrgService) CreateOrg(req OrgCreateRequest) (*Org, error) {
	// Check if organization with name already exists
	var existingOrg Org
	if err := s.db.Where("name = ?", req.Name).First(&existingOrg).Error; err == nil {
		return nil, errors.New("organization with this name already exists")
	}

	org := Org{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	result := s.db.Create(&org)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create organization: %w", result.Error)
	}

	return &org, nil
}

func (s *OrgService) UpdateOrg(id uuid.UUID, req OrgUpdateRequest) (*Org, error) {
	var org Org
	result := s.db.First(&org, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get organization: %w", result.Error)
	}

	// Update fields if provided
	if req.Name != "" {
		// Check if name is already taken by another organization
		var existingOrg Org
		if err := s.db.Where("name = ? AND id != ?", req.Name, id).First(&existingOrg).Error; err == nil {
			return nil, errors.New("organization name is already taken by another organization")
		}
		org.Name = req.Name
	}
	if req.DisplayName != "" {
		org.DisplayName = req.DisplayName
	}
	if req.Description != "" {
		org.Description = req.Description
	}

	org.UpdatedAt = time.Now()

	result = s.db.Save(&org)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to update organization: %w", result.Error)
	}

	return &org, nil
}

func (s *OrgService) DeleteOrg(id uuid.UUID) error {
	result := s.db.Delete(&Org{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete organization: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("organization not found")
	}
	return nil
}
