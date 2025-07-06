package controllers

import (
	"log/slog"
	"net/http"

	"idmapp-go/dto"
	"idmapp-go/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrgMemberController struct {
	service *services.OrgMemberService
	logger  *slog.Logger
}

func NewOrgMemberController(service *services.OrgMemberService) *OrgMemberController {
	return &OrgMemberController{
		service: service,
		logger:  slog.Default(),
	}
}

func (c *OrgMemberController) HandleMemberOperation(ctx *gin.Context) {
	var req dto.OrgMemberOpRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Op == 1 { // ADD
		orgMember, err := c.service.AddMember(req.OrgID, req.EntityID, req.Type)
		if err != nil {
			c.logger.Error("Failed to add org member", "error", err, "orgID", req.OrgID, "entityID", req.EntityID, "type", req.Type)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, orgMember)
		return
	} else if req.Op == 2 { // REMOVE
		removed, err := c.service.RemoveMember(req.OrgID, req.EntityID)
		if err != nil {
			c.logger.Error("Failed to remove org member", "error", err, "orgID", req.OrgID, "entityID", req.EntityID)
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

func (c *OrgMemberController) GetAllMembers(ctx *gin.Context) {
	members, err := c.service.GetAllMembers()
	if err != nil {
		c.logger.Error("Failed to get org members", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get org members"})
		return
	}

	ctx.JSON(http.StatusOK, members)
}

func (c *OrgMemberController) GetMembersByOrgID(ctx *gin.Context) {
	orgIDStr := ctx.Param("orgId")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	members, err := c.service.GetMembersByOrgID(orgID)
	if err != nil {
		c.logger.Error("Failed to get org members by org ID", "error", err, "orgID", orgID)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get org members"})
		return
	}

	ctx.JSON(http.StatusOK, members)
}

func (c *OrgMemberController) GetMembersByEntityID(ctx *gin.Context) {
	entityIDStr := ctx.Param("entityId")
	entityID, err := uuid.Parse(entityIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}

	members, err := c.service.GetMembersByEntityID(entityID)
	if err != nil {
		c.logger.Error("Failed to get org members by entity ID", "error", err, "entityID", entityID)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get org members"})
		return
	}

	ctx.JSON(http.StatusOK, members)
}
