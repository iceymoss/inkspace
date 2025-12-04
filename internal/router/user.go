package router

import (
	"mysite/internal/handler"
	"mysite/internal/middleware"

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
			public.GET("/articles/:id", articleHandler.GetDetail)
			public.GET("/articles/:id/is-liked", likeHandler.CheckArticleLiked)

			// Comments (public read)
			public.GET("/comments", commentHandler.GetList)
			public.GET("/comments/:id/is-liked", likeHandler.CheckCommentLiked)
			
			// Comment likes (public, can be done by guests)
			public.POST("/comments/:id/like", likeHandler.LikeComment)
			public.DELETE("/comments/:id/like", likeHandler.UnlikeComment)

			// Categories and Tags
			public.GET("/categories", categoryHandler.GetList)
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

			// Tags (users can create their own tags)
			protected.POST("/tags", tagHandler.Create)

			// Comments
			protected.POST("/comments", commentHandler.Create)
			protected.DELETE("/comments/:id", commentHandler.Delete)

			// Likes (articles require auth, comments are public)
			protected.POST("/articles/:id/like", likeHandler.LikeArticle)
			protected.DELETE("/articles/:id/like", likeHandler.UnlikeArticle)

			// Follow
			protected.POST("/users/:id/follow", followHandler.Follow)
			protected.DELETE("/users/:id/follow", followHandler.Unfollow)

			// Favorites
			protected.POST("/articles/:id/favorite", favoriteHandler.AddFavorite)
			protected.DELETE("/articles/:id/favorite", favoriteHandler.RemoveFavorite)
			protected.GET("/articles/:id/is-favorited", favoriteHandler.CheckFavorited)
			protected.GET("/favorites", favoriteHandler.GetMyFavorites)

			// Notifications
			protected.GET("/notifications", notificationHandler.GetList)
			protected.PUT("/notifications/:id/read", notificationHandler.MarkAsRead)
			protected.PUT("/notifications/read-all", notificationHandler.MarkAllAsRead)
			protected.GET("/notifications/unread-count", notificationHandler.GetUnreadCount)
			protected.DELETE("/notifications/:id", notificationHandler.Delete)
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

