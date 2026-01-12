package router

import (
	"github.com/iceymoss/inkspace/internal/handler"
	"github.com/iceymoss/inkspace/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupUserRouter 设置用户服务路由
func SetupUserRouter() *gin.Engine {
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
	adHandler := handler.NewAdHandler()

	// API routes
	api := r.Group("/api")
	{
		// Public routes
		public := api.Group("")
		{
			// Authentication
			public.POST("/register", userHandler.Register)
			public.POST("/login", userHandler.Login)

			// Articles (public read)
			public.GET("/articles", articleHandler.GetList)
			public.GET("/articles/recommended", articleHandler.GetRecommended)
			public.GET("/articles/hot", articleHandler.GetHotArticles)

			// 状态检查API（可选认证）
			publicWithOptionalAuth := api.Group("")
			publicWithOptionalAuth.Use(middleware.OptionalAuthMiddleware())
			{
				// 文章详情（需要可选认证，以便作者可以查看自己的私有/草稿文章）
				publicWithOptionalAuth.GET("/articles/:id", articleHandler.GetDetail)
				publicWithOptionalAuth.GET("/articles/:id/is-liked", likeHandler.CheckArticleLiked)
				publicWithOptionalAuth.GET("/articles/:id/is-favorited", favoriteHandler.CheckFavorited)
				// 作品详情（需要可选认证，以便作者可以查看自己的待审核/审核不通过的作品）
				publicWithOptionalAuth.GET("/works/:id", workHandler.GetDetail)
				publicWithOptionalAuth.GET("/works/:id/liked", likeHandler.CheckWorkLiked)
				publicWithOptionalAuth.GET("/works/:id/favorited", favoriteHandler.CheckWorkFavorited)
			}

			// Comments (public read)
			public.GET("/comments", commentHandler.GetList)
			public.GET("/comments/replies/:root_id", commentHandler.GetReplies) // 获取子评论分页列表
			public.GET("/comments/:id/is-liked", likeHandler.CheckCommentLiked)

			// Comment likes (public, can be done by guests)
			public.POST("/comments/:id/like", likeHandler.LikeComment)
			public.DELETE("/comments/:id/like", likeHandler.UnlikeComment)

			// Categories and Tags
			public.GET("/categories", categoryHandler.GetList)
			public.GET("/tags", tagHandler.GetList)

			// Works
			public.GET("/works", workHandler.GetList)
			public.GET("/works/recommended", workHandler.GetRecommended)
			public.GET("/works/hot", workHandler.GetHotWorks)
			// 作品详情已移到 publicWithOptionalAuth 组中

			// Work Comments (public read, same as article comments)
			// Comments endpoint handles both article and work comments via query params

			// User Profile (public)
			public.GET("/users/search", userHandler.SearchUsers)
			public.GET("/users/:id", userHandler.GetUserProfile)
			public.GET("/users/:id/articles", articleHandler.GetUserArticles)
			public.GET("/users/:id/works", workHandler.GetUserWorks)
		}

		// 可选认证的路由（支持未登录访问，但登录后会有额外信息）
		publicWithOptionalAuth := api.Group("")
		publicWithOptionalAuth.Use(middleware.OptionalAuthMiddleware())
		{
			// 关注统计（支持可选认证，以便显示当前用户的关注状态）
			publicWithOptionalAuth.GET("/users/:id/follow-stats", followHandler.GetFollowStats)

			// Links
			public.GET("/links", linkHandler.GetList)

			// Settings (public)
			public.GET("/settings/public", settingHandler.GetPublicSettings)

			// Ads (public)
			public.GET("/ads", adHandler.GetAdsByPositionCode)
			public.POST("/ads/:id/click", adHandler.RecordAdClick)
			public.POST("/ads/:id/view", adHandler.RecordAdView)
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
			protected.POST("/upload/photo", uploadHandler.UploadPhoto) // 摄影作品原图上传

			// Articles (author can manage their own articles)
			protected.GET("/articles/:id/edit", articleHandler.GetEdit) // 编辑页专用API，需要权限检查
			protected.POST("/articles", articleHandler.Create)
			protected.PUT("/articles/:id", articleHandler.Update)
			protected.DELETE("/articles/:id", articleHandler.Delete)
			// Markdown image upload (public, but rate limited)
			public.POST("/upload/markdown-image", uploadHandler.UploadMarkdownImage)

			// Works (author can manage their own works)
			protected.GET("/works/:id/edit", workHandler.GetEdit) // 编辑页专用API，需要权限检查
			protected.POST("/works", workHandler.Create)
			protected.PUT("/works/:id", workHandler.Update)
			protected.DELETE("/works/:id", workHandler.Delete)
			protected.GET("/works/quota", workHandler.GetQuotaUsage)
			protected.GET("/works/my", workHandler.GetMyWorks)

			// Tags (users can create their own tags)
			protected.POST("/tags", tagHandler.Create)

			// Comments
			protected.POST("/comments", commentHandler.Create)
			protected.DELETE("/comments/:id", commentHandler.Delete)

			// Likes (articles and works require auth, comments are public)
			protected.POST("/articles/:id/like", likeHandler.LikeArticle)
			protected.DELETE("/articles/:id/like", likeHandler.UnlikeArticle)
			protected.POST("/works/:id/like", likeHandler.LikeWork)

			// Follow
			protected.POST("/users/:id/follow", followHandler.Follow)
			protected.DELETE("/users/:id/follow", followHandler.Unfollow)
			// 用户关注/粉丝列表（仅本人可访问）
			protected.GET("/users/:id/following", followHandler.GetFollowingList)
			protected.GET("/users/:id/followers", followHandler.GetFollowerList)

			// Favorites
			protected.POST("/articles/:id/favorite", favoriteHandler.AddFavorite)
			protected.DELETE("/articles/:id/favorite", favoriteHandler.RemoveFavorite)
			protected.POST("/works/:id/favorite", favoriteHandler.AddWorkFavorite)
			protected.DELETE("/works/:id/favorite", favoriteHandler.RemoveWorkFavorite)
			protected.GET("/favorites", favoriteHandler.GetMyFavorites)
			// 用户收藏列表（仅本人可访问）
			protected.GET("/users/:id/favorites", favoriteHandler.GetUserFavorites)

			// Notifications
			protected.GET("/notifications", notificationHandler.GetNotifications)
			protected.GET("/notifications/unread-count", notificationHandler.GetUnreadCount)
			protected.PUT("/notifications/:id/read", notificationHandler.MarkAsRead)
			protected.PUT("/notifications/read-all", notificationHandler.MarkAllAsRead)
			protected.DELETE("/notifications/:id", notificationHandler.DeleteNotification)
			protected.DELETE("/notifications/read-all", notificationHandler.DeleteAllRead)
		}
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Serve static files (uploads)
	r.Static("/uploads", "./uploads")

	return r
}
