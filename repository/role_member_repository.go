package repository

import (
	"idmapp-go/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleMemberRepository struct {
	db *gorm.DB
}

func NewRoleMemberRepository(db *gorm.DB) *RoleMemberRepository {
	return &RoleMemberRepository{db: db}
}

func (r *RoleMemberRepository) FindByID(id uuid.UUID) (*models.RoleMember, error) {
	var roleMember models.RoleMember
	if err := r.db.First(&roleMember, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &roleMember, nil
}

func (r *RoleMemberRepository) FindByRoleID(roleID uuid.UUID) ([]models.RoleMember, error) {
	var members []models.RoleMember
	if err := r.db.Where("role_id = ?", roleID).Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

func (r *RoleMemberRepository) FindByEntityID(entityID uuid.UUID) ([]models.RoleMember, error) {
	var members []models.RoleMember
	if err := r.db.Where("entity_id = ?", entityID).Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

func (r *RoleMemberRepository) FindAll() ([]models.RoleMember, error) {
	var members []models.RoleMember
	if err := r.db.Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

func (r *RoleMemberRepository) Save(roleMember *models.RoleMember) error {
	return r.db.Create(roleMember).Error
}

func (r *RoleMemberRepository) DeleteByRoleIDAndEntityID(roleID, entityID uuid.UUID) (int64, error) {
	result := r.db.Where("role_id = ? AND entity_id = ?", roleID, entityID).Delete(&models.RoleMember{})
	return result.RowsAffected, result.Error
}
