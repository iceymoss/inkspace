package router

import (
	"github.com/iceymoss/inkspace/internal/handler"
	"github.com/iceymoss/inkspace/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupAdminRouter 设置管理后台路由
func SetupAdminRouter() *gin.Engine {
	r := gin.Default()

	// Middleware
	r.Use(middleware.CORSMiddleware())

	// Handlers
	userHandler := handler.NewUserHandler()
	articleHandler := handler.NewArticleHandler()
	commentHandler := handler.NewCommentHandler()
	categoryHandler := handler.NewCategoryHandler()
	tagHandler := handler.NewTagHandler()
	workHandler := handler.NewWorkHandler()
	linkHandler := handler.NewLinkHandler()
	settingHandler := handler.NewSettingHandler()
	adminAuthHandler := handler.NewAdminAuthHandler()
	uploadHandler := handler.NewUploadHandler()
	adHandler := handler.NewAdHandler()

	// 注意：管理后台需要完整的handler来处理查询和管理操作

	// API routes
	api := r.Group("/api")
	{
		// Admin authentication routes (public)
		adminAuth := api.Group("/admin/auth")
		{
			adminAuth.POST("/login", adminAuthHandler.AdminLogin)
		}

		// Admin authenticated routes
		adminAuthRoutes := api.Group("/admin/auth")
		adminAuthRoutes.Use(middleware.AdminAuthMiddleware())
		{
			adminAuthRoutes.GET("/profile", adminAuthHandler.GetAdminProfile)
			adminAuthRoutes.POST("/logout", adminAuthHandler.AdminLogout)
		}

		// Upload routes (authenticated)
		upload := api.Group("/upload")
		upload.Use(middleware.AdminAuthMiddleware())
		{
			upload.POST("/image", uploadHandler.UploadImage)
			upload.POST("/avatar", uploadHandler.UploadAvatar)
		}

		// Admin routes (require admin authentication)
		admin := api.Group("/admin")
		admin.Use(middleware.AdminAuthMiddleware())
		{
			// Users management
			admin.GET("/users", userHandler.GetUserList)
			admin.PUT("/users/:id/status", userHandler.UpdateUserStatus)
			admin.PUT("/users/:id/role", userHandler.UpdateUserRole)
			admin.DELETE("/users/:id", userHandler.DeleteUser)

			// Articles management
			admin.GET("/articles", articleHandler.GetList)
			admin.GET("/articles/:id", articleHandler.GetDetail)
			admin.POST("/articles", articleHandler.Create)
			admin.PUT("/articles/:id", articleHandler.Update)
			admin.PUT("/articles/:id/recommend", articleHandler.SetRecommend)
			admin.DELETE("/articles/:id", articleHandler.Delete)

			// Works management
			admin.GET("/works", workHandler.GetList)
			admin.GET("/works/:id", workHandler.GetDetail)
			admin.POST("/works", workHandler.Create)
			admin.PUT("/works/:id", workHandler.Update)
			admin.PUT("/works/:id/recommend", workHandler.SetRecommend)
			admin.PUT("/works/:id/status", workHandler.UpdateWorkStatus)
			admin.DELETE("/works/:id", workHandler.Delete)

			// Categories management
			admin.GET("/categories", categoryHandler.GetList)
			admin.POST("/categories", categoryHandler.Create)
			admin.PUT("/categories/:id", categoryHandler.Update)
			admin.DELETE("/categories/:id", categoryHandler.Delete)

			// Tags management
			admin.GET("/tags", tagHandler.GetList)
			admin.POST("/tags", tagHandler.Create)
			admin.PUT("/tags/:id", tagHandler.Update)
			admin.DELETE("/tags/:id", tagHandler.Delete)

			// Comments management
			admin.GET("/comments", commentHandler.GetList)
			admin.PUT("/comments/:id/status", commentHandler.UpdateStatus)
			admin.DELETE("/comments/:id", commentHandler.Delete)

			// Links management
			admin.GET("/links", linkHandler.GetList)
			admin.POST("/links", linkHandler.Create)
			admin.PUT("/links/:id", linkHandler.Update)
			admin.DELETE("/links/:id", linkHandler.Delete)

			// Settings management
			admin.GET("/settings", settingHandler.GetAllSettings)
			admin.PUT("/settings", settingHandler.UpdateSetting)
			admin.PUT("/settings/batch", settingHandler.BatchUpdateSettings)
			admin.DELETE("/settings/:key", settingHandler.DeleteSetting)

			// Ad Positions management
			admin.GET("/ad-positions", adHandler.GetPositionList)
			admin.GET("/ad-positions/:id", adHandler.GetPositionByID)
			admin.POST("/ad-positions", adHandler.CreatePosition)
			admin.PUT("/ad-positions/:id", adHandler.UpdatePosition)
			admin.DELETE("/ad-positions/:id", adHandler.DeletePosition)

			// Advertisements management
			admin.GET("/advertisements", adHandler.GetAdvertisementList)
			admin.GET("/advertisements/:id", adHandler.GetAdvertisementByID)
			admin.POST("/advertisements", adHandler.CreateAdvertisement)
			admin.PUT("/advertisements/:id", adHandler.UpdateAdvertisement)
			admin.DELETE("/advertisements/:id", adHandler.DeleteAdvertisement)

			// Ad Placements management
			admin.GET("/ad-placements", adHandler.GetPlacementList)
			admin.GET("/ad-placements/:id", adHandler.GetPlacementByID)
			admin.POST("/ad-placements", adHandler.CreatePlacement)
			admin.PUT("/ad-placements/:id", adHandler.UpdatePlacement)
			admin.DELETE("/ad-placements/:id", adHandler.DeletePlacement)
		}
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "admin"})
	})

	// Serve static files (uploads)
	r.Static("/uploads", "./uploads")

	return r
}
