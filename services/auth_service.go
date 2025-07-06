package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/openfga/go-sdk/client"
	"github.com/openfga/go-sdk/credentials"
)

type AuthorizationService struct {
	fgaClient *client.OpenFgaClient
	logger    *slog.Logger
}

func NewAuthorizationService(apiURL, storeID, apiToken string) (*AuthorizationService, error) {
	fgaClient, err := client.NewSdkClient(&client.ClientConfiguration{
		ApiUrl:  apiURL,
		StoreId: storeID,
		Credentials: &credentials.Credentials{
			Method: credentials.CredentialsMethodApiToken,
			Config: &credentials.Config{
				ApiToken: apiToken,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenFGA client: %w", err)
	}

	return &AuthorizationService{
		fgaClient: fgaClient,
		logger:    slog.Default(),
	}, nil
}

func (s *AuthorizationService) CheckAccess(userID, relation, resourceID string) bool {
	body := client.ClientCheckRequest{
		User:     userID,
		Relation: relation,
		Object:   resourceID,
	}

	resp, err := s.fgaClient.Check(context.Background()).Body(body).Execute()
	if err != nil {
		s.logger.Error("Error checking permission", "user", userID, "relation", relation, "resource", resourceID, "error", err)
		return false
	}

	allowed := resp.GetAllowed()
	s.logger.Debug("Permission check", "user", userID, "relation", relation, "resource", resourceID, "result", map[bool]string{true: "ALLOWED", false: "DENIED"}[allowed])
	return allowed
}

func (s *AuthorizationService) GrantPermission(userID, relation, resourceID string) bool {
	writes := []client.ClientTupleKey{{
		User:     userID,
		Relation: relation,
		Object:   resourceID,
	}}
	body := client.ClientWriteRequest{
		Writes: writes,
	}

	_, err := s.fgaClient.Write(context.Background()).Body(body).Execute()
	if err != nil {
		s.logger.Error("Error granting permission", "user", userID, "relation", relation, "resource", resourceID, "error", err)
		return false
	}

	s.logger.Info("Granted permission", "user", userID, "relation", relation, "resource", resourceID)
	return true
}

func (s *AuthorizationService) RevokePermission(userID, relation, resourceID string) bool {
	deletes := []client.ClientTupleKeyWithoutCondition{{
		User:     userID,
		Relation: relation,
		Object:   resourceID,
	}}
	body := client.ClientWriteRequest{
		Deletes: deletes,
	}

	_, err := s.fgaClient.Write(context.Background()).Body(body).Execute()
	if err != nil {
		s.logger.Error("Error revoking permission", "user", userID, "relation", relation, "resource", resourceID, "error", err)
		return false
	}

	s.logger.Info("Revoked permission", "user", userID, "relation", relation, "resource", resourceID)
	return true
}

// Convenience methods for common authorization checks
func (s *AuthorizationService) CanManageUser(userID, targetUserID string) bool {
	return s.CheckAccess(userID, "manage", "user:"+targetUserID)
}

func (s *AuthorizationService) CanManageGroup(userID, groupID string) bool {
	return s.CheckAccess(userID, "manage", "group:"+groupID)
}

func (s *AuthorizationService) CanManageRole(userID, roleID string) bool {
	return s.CheckAccess(userID, "manage", "role:"+roleID)
}

func (s *AuthorizationService) CanManageOrg(userID, orgID string) bool {
	return s.CheckAccess(userID, "manage", "org:"+orgID)
} 