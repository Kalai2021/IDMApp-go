package services

import (
	"idmapp-go/models"
	"idmapp-go/repository"
	"github.com/google/uuid"
)

type UserGroupMemberService struct {
	repo *repository.UserGroupMemberRepository
}

func NewUserGroupMemberService(repo *repository.UserGroupMemberRepository) *UserGroupMemberService {
	return &UserGroupMemberService{repo: repo}
}

func (s *UserGroupMemberService) AddMember(groupId, userId uuid.UUID) (*models.UserGroupMember, error) {
	member := &models.UserGroupMember{
		GroupID:   groupId,
		UserID:    userId,
	}
	if err := s.repo.Save(member); err != nil {
		return nil, err
	}
	return member, nil
}

func (s *UserGroupMemberService) RemoveMember(groupId, userId uuid.UUID) (bool, error) {
	affected, err := s.repo.DeleteByGroupIDAndUserID(groupId, userId)
	if err != nil {
		return false, err
	}
	return affected > 0, nil
} 