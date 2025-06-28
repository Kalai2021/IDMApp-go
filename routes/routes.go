package routes

import (
	"idmapp-go/controllers"
	"idmapp-go/database"
	"idmapp-go/internal/group"
	"idmapp-go/internal/member"
	"idmapp-go/internal/org"
	"idmapp-go/internal/role"
	"idmapp-go/internal/user"
	"idmapp-go/middleware"
	"idmapp-go/repository"
	"idmapp-go/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// Serve static files
	router.Static("/static", "./templates")
	router.StaticFile("/test.html", "./test.html")

	// Initialize services
	userService := user.NewUserService(database.GetDB())
	groupService := group.NewGroupService(database.GetDB())
	roleService := role.NewRoleService(database.GetDB())
	orgService := org.NewOrgService(database.GetDB())
	memberService := member.NewMemberService(database.GetDB())
	pkceService := services.NewPKCEService(database.GetDB())

	// Initialize repositories for member services
	db := database.GetDB()
	orgMemberRepo := repository.NewOrgMemberRepository(db)
	roleMemberRepo := repository.NewRoleMemberRepository(db)

	// Initialize member services
	orgMemberService := services.NewOrgMemberService(orgMemberRepo)
	roleMemberService := services.NewRoleMemberService(roleMemberRepo)

	// Initialize controllers
	userController := user.NewUserController(userService, pkceService)
	groupController := group.NewGroupController(groupService)
	roleController := role.NewRoleController(roleService)
	orgController := org.NewOrgController(orgService)
	memberController := member.NewMemberController(memberService)
	orgMemberController := controllers.NewOrgMemberController(orgMemberService)
	roleMemberController := controllers.NewRoleMemberController(roleMemberService)
	pkceController := controllers.NewPKCEController(pkceService, userService)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes (no authentication required)
		v1.POST("/auth/login", userController.Login)

		// PKCE Authentication routes (public)
		pkce := v1.Group("/auth/pkce")
		{
			pkce.GET("/config", pkceController.GetPKCEConfig)
			pkce.GET("/jwks", pkceController.GetJWKS)
			pkce.GET("/authorize", pkceController.InitiatePKCEAuthGET)
			pkce.POST("/token", pkceController.ExchangeCodeForToken)
			pkce.POST("/refresh", pkceController.RefreshToken)
		}

		// Login form routes (public)
		router.GET("/login", controllers.ShowLoginForm)
		router.POST("/login", controllers.HandleLogin)
		router.GET("/logout", controllers.Logout)

		// Protected routes (authentication required)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("", userController.GetAllUsers)
				users.GET("/:id", userController.GetUser)
				users.POST("", userController.CreateUser)
				users.PUT("/:id", userController.UpdateUser)
				users.DELETE("/:id", userController.DeleteUser)
			}

			// Group routes
			groups := protected.Group("/groups")
			{
				groups.GET("", groupController.GetAllGroups)
				groups.GET("/:id", groupController.GetGroup)
				groups.POST("", groupController.CreateGroup)
				groups.PUT("/:id", groupController.UpdateGroup)
				groups.DELETE("/:id", groupController.DeleteGroup)
			}

			// Role routes
			roles := protected.Group("/roles")
			{
				roles.GET("", roleController.GetAllRoles)
				roles.GET("/:id", roleController.GetRole)
				roles.POST("", roleController.CreateRole)
				roles.PUT("/:id", roleController.UpdateRole)
				roles.DELETE("/:id", roleController.DeleteRole)
			}

			// Organization routes
			orgs := protected.Group("/orgs")
			{
				orgs.GET("", orgController.GetAllOrgs)
				orgs.GET("/:id", orgController.GetOrg)
				orgs.POST("", orgController.CreateOrg)
				orgs.PUT("/:id", orgController.UpdateOrg)
				orgs.DELETE("/:id", orgController.DeleteOrg)
			}

			// Member routes (User-Group management)
			members := protected.Group("/groupmembers")
			{
				members.GET("", memberController.GetAllMembers)
				members.GET("/group/:groupId", memberController.GetMembersByGroupID)
				members.GET("/user/:userId", memberController.GetMembersByUserID)
				members.POST("", memberController.AddMember)
			}

			// Organization Member routes
			orgMembers := protected.Group("/orgmembers")
			{
				orgMembers.GET("", orgMemberController.GetAllMembers)
				orgMembers.GET("/org/:orgId", orgMemberController.GetMembersByOrgID)
				orgMembers.POST("", orgMemberController.HandleMemberOperation)
			}

			// Role Member routes
			roleMembers := protected.Group("/rolemembers")
			{
				roleMembers.GET("", roleMemberController.GetAllMembers)
				roleMembers.GET("/role/:roleId", roleMemberController.GetMembersByRoleID)
				roleMembers.GET("/entity/:entityId", roleMemberController.GetMembersByEntityID)
				roleMembers.POST("", roleMemberController.HandleMemberOperation)
			}
		}
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "IDM App is running",
		})
	})

	// OIDC Discovery endpoint (well-known)
	router.GET("/.well-known/openid-configuration", pkceController.GetOIDCConfig)
}
