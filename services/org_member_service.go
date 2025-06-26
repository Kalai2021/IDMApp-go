package services

import (
	"idmapp-go/models"
	"idmapp-go/repository"

	"github.com/google/uuid"
)

type OrgMemberService struct {
	repo *repository.OrgMemberRepository
}

func NewOrgMemberService(repo *repository.OrgMemberRepository) *OrgMemberService {
	return &OrgMemberService{repo: repo}
}

func (s *OrgMemberService) AddMember(orgId, entityId uuid.UUID, memberType string) (*models.OrgMember, error) {
	orgMember := &models.OrgMember{
		OrgID:    orgId,
		EntityID: entityId,
		Type:     memberType,
	}
	if err := s.repo.Save(orgMember); err != nil {
		return nil, err
	}
	return orgMember, nil
}

func (s *OrgMemberService) RemoveMember(orgId, entityId uuid.UUID) (bool, error) {
	affected, err := s.repo.DeleteByOrgIDAndEntityID(orgId, entityId)
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

func (s *OrgMemberService) GetMembersByOrgID(orgId uuid.UUID) ([]models.OrgMember, error) {
	return s.repo.FindByOrgID(orgId)
}

func (s *OrgMemberService) GetMembersByEntityID(entityId uuid.UUID) ([]models.OrgMember, error) {
	return s.repo.FindByEntityID(entityId)
}

func (s *OrgMemberService) GetAllMembers() ([]models.OrgMember, error) {
	return s.repo.FindAll()
}
