package repository

import (
	"idmapp-go/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrgMemberRepository struct {
	db *gorm.DB
}

func NewOrgMemberRepository(db *gorm.DB) *OrgMemberRepository {
	return &OrgMemberRepository{db: db}
}

func (r *OrgMemberRepository) FindByID(id uuid.UUID) (*models.OrgMember, error) {
	var orgMember models.OrgMember
	if err := r.db.First(&orgMember, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &orgMember, nil
}

func (r *OrgMemberRepository) FindByOrgID(orgID uuid.UUID) ([]models.OrgMember, error) {
	var members []models.OrgMember
	if err := r.db.Where("org_id = ?", orgID).Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

func (r *OrgMemberRepository) FindByEntityID(entityID uuid.UUID) ([]models.OrgMember, error) {
	var members []models.OrgMember
	if err := r.db.Where("entity_id = ?", entityID).Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

func (r *OrgMemberRepository) FindAll() ([]models.OrgMember, error) {
	var members []models.OrgMember
	if err := r.db.Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

func (r *OrgMemberRepository) Save(orgMember *models.OrgMember) error {
	return r.db.Create(orgMember).Error
}

func (r *OrgMemberRepository) DeleteByOrgIDAndEntityID(orgID, entityID uuid.UUID) (int64, error) {
	result := r.db.Where("org_id = ? AND entity_id = ?", orgID, entityID).Delete(&models.OrgMember{})
	return result.RowsAffected, result.Error
}
