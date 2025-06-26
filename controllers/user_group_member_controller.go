package controllers

import (
	"net/http"

	"idmapp-go/dto"
	"idmapp-go/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserGroupMemberController struct {
	service *services.UserGroupMemberService
	logger  *logrus.Logger
}

func NewUserGroupMemberController(service *services.UserGroupMemberService) *UserGroupMemberController {
	return &UserGroupMemberController{
		service: service,
		logger:  logrus.New(),
	}
}

func (c *UserGroupMemberController) HandleMemberOperation(ctx *gin.Context) {
	var req dto.UserGroupMemberOpRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Op == 1 { // ADD
		member, err := c.service.AddMember(req.GroupID, req.UserID)
		if err != nil {
			c.logger.Errorf("Failed to add user group member: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, member)
		return
	} else if req.Op == 2 { // REMOVE
		removed, err := c.service.RemoveMember(req.GroupID, req.UserID)
		if err != nil {
			c.logger.Errorf("Failed to remove user group member: %v", err)
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