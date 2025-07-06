package role

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RoleController struct {
	roleService *RoleService
	logger      *slog.Logger
}

func NewRoleController(roleService *RoleService) *RoleController {
	return &RoleController{
		roleService: roleService,
		logger:      slog.Default(),
	}
}

func (c *RoleController) GetAllRoles(ctx *gin.Context) {
	roles, err := c.roleService.GetAllRoles()
	if err != nil {
		c.logger.Error("Failed to get roles", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get roles"})
		return
	}

	ctx.JSON(http.StatusOK, roles)
}

func (c *RoleController) GetRole(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	role, err := c.roleService.GetRole(id)
	if err != nil {
		c.logger.Error("Failed to get role", "error", err, "roleID", id)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get role"})
		return
	}

	if role == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	ctx.JSON(http.StatusOK, role)
}

func (c *RoleController) CreateRole(ctx *gin.Context) {
	var req RoleCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role, err := c.roleService.CreateRole(req)
	if err != nil {
		c.logger.Error("Failed to create role", "error", err, "name", req.Name)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, role)
}

func (c *RoleController) UpdateRole(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	var req RoleUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role, err := c.roleService.UpdateRole(id, req)
	if err != nil {
		c.logger.Error("Failed to update role", "error", err, "roleID", id)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if role == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	ctx.JSON(http.StatusOK, role)
}

func (c *RoleController) DeleteRole(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	err = c.roleService.DeleteRole(id)
	if err != nil {
		c.logger.Error("Failed to delete role", "error", err, "roleID", id)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
