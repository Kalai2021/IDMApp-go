package controllers

import (
	"net/http"

	"idmapp-go/dto"
	"idmapp-go/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type RoleMemberController struct {
	service *services.RoleMemberService
	logger  *logrus.Logger
}

func NewRoleMemberController(service *services.RoleMemberService) *RoleMemberController {
	return &RoleMemberController{
		service: service,
		logger:  logrus.New(),
	}
}

func (c *RoleMemberController) HandleMemberOperation(ctx *gin.Context) {
	var req dto.RoleMemberOpRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Op == 1 { // ADD
		roleMember, err := c.service.AddMember(req.RoleID, req.EntityID, req.Type)
		if err != nil {
			c.logger.Errorf("Failed to add role member: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, roleMember)
		return
	} else if req.Op == 2 { // REMOVE
		removed, err := c.service.RemoveMember(req.RoleID, req.EntityID)
		if err != nil {
			c.logger.Errorf("Failed to remove role member: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if removed {
			ctx.Status(http.StatusOK)
		} else {
			ctx.Status(http.StatusNotFound)
		}
		return
	}
	ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid operation code"})
}

func (c *RoleMemberController) GetAllMembers(ctx *gin.Context) {
	members, err := c.service.GetAllMembers()
	if err != nil {
		c.logger.Errorf("Failed to get role members: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get role members"})
		return
	}

	ctx.JSON(http.StatusOK, members)
}

func (c *RoleMemberController) GetMembersByRoleID(ctx *gin.Context) {
	roleIDStr := ctx.Param("roleId")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	members, err := c.service.GetMembersByRoleID(roleID)
	if err != nil {
		c.logger.Errorf("Failed to get role members by role ID: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get role members"})
		return
	}

	ctx.JSON(http.StatusOK, members)
}

func (c *RoleMemberController) GetMembersByEntityID(ctx *gin.Context) {
	entityIDStr := ctx.Param("entityId")
	entityID, err := uuid.Parse(entityIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}

	members, err := c.service.GetMembersByEntityID(entityID)
	if err != nil {
		c.logger.Errorf("Failed to get role members by entity ID: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get role members"})
		return
	}

	ctx.JSON(http.StatusOK, members)
}
