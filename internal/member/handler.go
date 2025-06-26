package member

import (
	"net/http"

	"idmapp-go/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type MemberController struct {
	memberService *MemberService
	logger        *logrus.Logger
}

func NewMemberController(memberService *MemberService) *MemberController {
	return &MemberController{
		memberService: memberService,
		logger:        logrus.New(),
	}
}

func (c *MemberController) AddMember(ctx *gin.Context) {
	var req dto.MemberOpRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	member, err := c.memberService.ProcessMemberOperation(req)
	if err != nil {
		c.logger.Errorf("Failed to process member operation: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if member == nil {
		// Operation was successful but no member returned (e.g., removal)
		ctx.Status(http.StatusNoContent)
		return
	}

	ctx.JSON(http.StatusCreated, member)
}

func (c *MemberController) GetAllMembers(ctx *gin.Context) {
	members, err := c.memberService.GetAllMembers()
	if err != nil {
		c.logger.Errorf("Failed to get members: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get members"})
		return
	}

	ctx.JSON(http.StatusOK, members)
}

func (c *MemberController) GetMembersByGroupID(ctx *gin.Context) {
	groupIDStr := ctx.Param("groupId")
	groupID, err := uuid.Parse(groupIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	members, err := c.memberService.GetMembersByGroupID(groupID)
	if err != nil {
		c.logger.Errorf("Failed to get members by group ID: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get members"})
		return
	}

	ctx.JSON(http.StatusOK, members)
}

func (c *MemberController) GetMembersByUserID(ctx *gin.Context) {
	userIDStr := ctx.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	members, err := c.memberService.GetMembersByUserID(userID)
	if err != nil {
		c.logger.Errorf("Failed to get members by user ID: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get members"})
		return
	}

	ctx.JSON(http.StatusOK, members)
}
