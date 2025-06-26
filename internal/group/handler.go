package group

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type GroupController struct {
	groupService *GroupService
	logger       *logrus.Logger
}

func NewGroupController(groupService *GroupService) *GroupController {
	return &GroupController{
		groupService: groupService,
		logger:       logrus.New(),
	}
}

func (c *GroupController) GetAllGroups(ctx *gin.Context) {
	groups, err := c.groupService.GetAllGroups()
	if err != nil {
		c.logger.Errorf("Failed to get groups: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get groups"})
		return
	}

	ctx.JSON(http.StatusOK, groups)
}

func (c *GroupController) GetGroup(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	group, err := c.groupService.GetGroup(id)
	if err != nil {
		c.logger.Errorf("Failed to get group: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get group"})
		return
	}

	if group == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	ctx.JSON(http.StatusOK, group)
}

func (c *GroupController) CreateGroup(ctx *gin.Context) {
	var req GroupCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	group, err := c.groupService.CreateGroup(req)
	if err != nil {
		c.logger.Errorf("Failed to create group: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, group)
}

func (c *GroupController) UpdateGroup(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	var req GroupUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	group, err := c.groupService.UpdateGroup(id, req)
	if err != nil {
		c.logger.Errorf("Failed to update group: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if group == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	ctx.JSON(http.StatusOK, group)
}

func (c *GroupController) DeleteGroup(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	err = c.groupService.DeleteGroup(id)
	if err != nil {
		c.logger.Errorf("Failed to delete group: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
