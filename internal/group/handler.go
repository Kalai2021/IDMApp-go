package group

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GroupController struct {
	groupService *GroupService
	logger       *slog.Logger
}

func NewGroupController(groupService *GroupService) *GroupController {
	return &GroupController{
		groupService: groupService,
		logger:       slog.Default(),
	}
}

func (c *GroupController) GetAllGroups(ctx *gin.Context) {
	groups, err := c.groupService.GetAllGroups()
	if err != nil {
		c.logger.Error("Failed to get groups", "error", err)
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
		c.logger.Error("Failed to get group", "error", err, "groupID", id)
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
		c.logger.Error("Failed to create group", "error", err, "name", req.Name)
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
		c.logger.Error("Failed to update group", "error", err, "groupID", id)
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
		c.logger.Error("Failed to delete group", "error", err, "groupID", id)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
