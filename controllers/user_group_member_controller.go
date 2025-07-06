package controllers

import (
	"log/slog"
	"net/http"

	"idmapp-go/dto"
	"idmapp-go/services"

	"github.com/gin-gonic/gin"
)

type UserGroupMemberController struct {
	service *services.UserGroupMemberService
	logger  *slog.Logger
}

func NewUserGroupMemberController(service *services.UserGroupMemberService) *UserGroupMemberController {
	return &UserGroupMemberController{
		service: service,
		logger:  slog.Default(),
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
			c.logger.Error("Failed to add user group member", "error", err, "groupID", req.GroupID, "userID", req.UserID)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, member)
		return
	} else if req.Op == 2 { // REMOVE
		removed, err := c.service.RemoveMember(req.GroupID, req.UserID)
		if err != nil {
			c.logger.Error("Failed to remove user group member", "error", err, "groupID", req.GroupID, "userID", req.UserID)
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
