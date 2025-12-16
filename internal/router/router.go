package router

import (
	"github.com/iceymoss/inkspace/internal/handler"
	"github.com/iceymoss/inkspace/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter 旧版路由函数（已废弃，请使用 SetupUserRouter 或 SetupAdminRouter）
// 保留此函数仅为向后兼容
func SetupRouter() *gin.Engine {
	// 默认返回用户路由
	return SetupUserRouter()
}

func SetupRouterOld() *gin.Engine {
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
	followHandler := handler.NewFollowHandler()
	favoriteHandler := handler.NewFavoriteHandler()
	likeHandler := handler.NewLikeHandler()
	linkHandler := handler.NewLinkHandler()
	settingHandler := handler.NewSettingHandler()
	notificationHandler := handler.NewNotificationHandler()
	uploadHandler := handler.NewUploadHandler()

	// API routes
	api := r.Group("/api")
	{
		// Public routes
		public := api.Group("")
		{
			// Auth
			public.POST("/register", userHandler.Register)
			public.POST("/login", userHandler.Login)

			// Articles
			public.GET("/articles", articleHandler.GetList)
			public.GET("/articles/:id", articleHandler.GetDetail)
			public.POST("/articles/:id/like", likeHandler.LikeArticle)
			public.DELETE("/articles/:id/like", likeHandler.UnlikeArticle)
			public.GET("/articles/:id/is-liked", likeHandler.CheckArticleLiked)

			// Comments
			public.GET("/comments", commentHandler.GetList)
			public.POST("/comments", commentHandler.Create)
			public.POST("/comments/:id/like", likeHandler.LikeComment)
			public.DELETE("/comments/:id/like", likeHandler.UnlikeComment)

			// Categories
			public.GET("/categories", categoryHandler.GetList)

			// Tags
			public.GET("/tags", tagHandler.GetList)

			// Works
			public.GET("/works", workHandler.GetList)
			public.GET("/works/:id", workHandler.GetDetail)

			// User Profile (public)
			public.GET("/users/:id", userHandler.GetUserProfile)
			public.GET("/users/:id/articles", articleHandler.GetUserArticles)
			public.GET("/users/:id/follow-stats", followHandler.GetFollowStats)
			public.GET("/users/:id/following", followHandler.GetFollowingList)
			public.GET("/users/:id/followers", followHandler.GetFollowerList)
			public.GET("/users/:id/favorites", favoriteHandler.GetUserFavorites)

			// Links
			public.GET("/links", linkHandler.GetList)

			// Settings (public)
			public.GET("/settings/public", settingHandler.GetPublicSettings)
		}

		// Protected routes (require authentication)
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// User
			protected.GET("/profile", userHandler.GetProfile)
			protected.PUT("/profile", userHandler.UpdateProfile)
			protected.PUT("/profile/password", userHandler.ChangePassword)

			// Upload
			protected.POST("/upload/image", uploadHandler.UploadImage)
			protected.POST("/upload/avatar", uploadHandler.UploadAvatar)

			// Articles (author can manage their own articles)
			protected.POST("/articles", articleHandler.Create)
			protected.PUT("/articles/:id", articleHandler.Update)
			protected.DELETE("/articles/:id", articleHandler.Delete)

			// Comments (author can delete their own comments)
			protected.DELETE("/comments/:id", commentHandler.Delete)

			// Follow
			protected.POST("/users/:id/follow", followHandler.Follow)
			protected.DELETE("/users/:id/follow", followHandler.Unfollow)

			// Favorite
			protected.POST("/articles/:id/favorite", favoriteHandler.AddFavorite)
			protected.DELETE("/articles/:id/favorite", favoriteHandler.RemoveFavorite)
			protected.GET("/articles/:id/is-favorited", favoriteHandler.CheckFavorited)
			protected.GET("/favorites", favoriteHandler.GetMyFavorites)

			// Notifications
			protected.GET("/notifications", notificationHandler.GetNotifications)
			protected.PUT("/notifications/:id/read", notificationHandler.MarkAsRead)
			protected.PUT("/notifications/read-all", notificationHandler.MarkAllAsRead)
			protected.GET("/notifications/unread-count", notificationHandler.GetUnreadCount)
			protected.DELETE("/notifications/:id", notificationHandler.DeleteNotification)
		}

		// ⚠️ 注意：所有管理后台路由已移至独立服务
		// 用户服务不再提供 /api/admin/* 路由
		// 请使用独立的管理后台服务 (端口 8083)
		// 启动命令: go run cmd/admin/main.go
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Serve static files (uploads)
	r.Static("/uploads", "./uploads")

	return r
}
