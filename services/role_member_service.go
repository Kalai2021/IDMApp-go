package services

import (
	"idmapp-go/models"
	"idmapp-go/repository"

	"github.com/google/uuid"
)

type RoleMemberService struct {
	repo *repository.RoleMemberRepository
}

func NewRoleMemberService(repo *repository.RoleMemberRepository) *RoleMemberService {
	return &RoleMemberService{repo: repo}
}

func (s *RoleMemberService) AddMember(roleId, entityId uuid.UUID, memberType string) (*models.RoleMember, error) {
	roleMember := &models.RoleMember{
		RoleID:   roleId,
		EntityID: entityId,
		Type:     memberType,
	}
	if err := s.repo.Save(roleMember); err != nil {
		return nil, err
	}
	return roleMember, nil
}

func (s *RoleMemberService) RemoveMember(roleId, entityId uuid.UUID) (bool, error) {
	affected, err := s.repo.DeleteByRoleIDAndEntityID(roleId, entityId)
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

func (s *RoleMemberService) GetMembersByRoleID(roleId uuid.UUID) ([]models.RoleMember, error) {
	return s.repo.FindByRoleID(roleId)
}

func (s *RoleMemberService) GetMembersByEntityID(entityId uuid.UUID) ([]models.RoleMember, error) {
	return s.repo.FindByEntityID(entityId)
}

func (s *RoleMemberService) GetAllMembers() ([]models.RoleMember, error) {
	return s.repo.FindAll()
}
