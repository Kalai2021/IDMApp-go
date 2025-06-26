package org

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type OrgController struct {
	orgService *OrgService
	logger     *logrus.Logger
}

func NewOrgController(orgService *OrgService) *OrgController {
	return &OrgController{
		orgService: orgService,
		logger:     logrus.New(),
	}
}

func (c *OrgController) GetAllOrgs(ctx *gin.Context) {
	orgs, err := c.orgService.GetAllOrgs()
	if err != nil {
		c.logger.Errorf("Failed to get organizations: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get organizations"})
		return
	}

	ctx.JSON(http.StatusOK, orgs)
}

func (c *OrgController) GetOrg(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	org, err := c.orgService.GetOrg(id)
	if err != nil {
		c.logger.Errorf("Failed to get organization: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get organization"})
		return
	}

	if org == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}

	ctx.JSON(http.StatusOK, org)
}

func (c *OrgController) CreateOrg(ctx *gin.Context) {
	var req OrgCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	org, err := c.orgService.CreateOrg(req)
	if err != nil {
		c.logger.Errorf("Failed to create organization: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, org)
}

func (c *OrgController) UpdateOrg(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	var req OrgUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	org, err := c.orgService.UpdateOrg(id, req)
	if err != nil {
		c.logger.Errorf("Failed to update organization: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if org == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}

	ctx.JSON(http.StatusOK, org)
}

func (c *OrgController) DeleteOrg(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	err = c.orgService.DeleteOrg(id)
	if err != nil {
		c.logger.Errorf("Failed to delete organization: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
