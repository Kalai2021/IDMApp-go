package repository

import (
	"idmapp-go/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserGroupMemberRepository struct {
	db *gorm.DB
}

func NewUserGroupMemberRepository(db *gorm.DB) *UserGroupMemberRepository {
	return &UserGroupMemberRepository{db: db}
}

func (r *UserGroupMemberRepository) FindByID(id uuid.UUID) (*models.UserGroupMember, error) {
	var member models.UserGroupMember
	if err := r.db.First(&member, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *UserGroupMemberRepository) FindByGroupID(groupID uuid.UUID) ([]models.UserGroupMember, error) {
	var members []models.UserGroupMember
	if err := r.db.Where("group_id = ?", groupID).Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

func (r *UserGroupMemberRepository) Save(member *models.UserGroupMember) error {
	return r.db.Create(member).Error
}

func (r *UserGroupMemberRepository) DeleteByGroupIDAndUserID(groupID, userID uuid.UUID) (int64, error) {
	result := r.db.Where("group_id = ? AND user_id = ?", groupID, userID).Delete(&models.UserGroupMember{})
	return result.RowsAffected, result.Error
} 