package router

import (
	"mysite/internal/handler"
	"mysite/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
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
			public.POST("/articles/:id/like", articleHandler.Like)

			// Comments
			public.GET("/comments", commentHandler.GetList)
			public.POST("/comments", commentHandler.Create)

			// Categories
			public.GET("/categories", categoryHandler.GetList)

			// Tags
			public.GET("/tags", tagHandler.GetList)

			// Works
			public.GET("/works", workHandler.GetList)
			public.GET("/works/:id", workHandler.GetDetail)
		}

		// Protected routes (require authentication)
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// User
			protected.GET("/profile", userHandler.GetProfile)
			protected.PUT("/profile", userHandler.UpdateProfile)

			// Articles (author can manage their own articles)
			protected.POST("/articles", articleHandler.Create)
			protected.PUT("/articles/:id", articleHandler.Update)
			protected.DELETE("/articles/:id", articleHandler.Delete)

			// Comments (author can delete their own comments)
			protected.DELETE("/comments/:id", commentHandler.Delete)
		}

		// Admin routes
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
		{
			// Users management
			admin.GET("/users", userHandler.GetUserList)

			// Categories management
			admin.POST("/categories", categoryHandler.Create)
			admin.PUT("/categories/:id", categoryHandler.Update)
			admin.DELETE("/categories/:id", categoryHandler.Delete)

			// Tags management
			admin.POST("/tags", tagHandler.Create)
			admin.PUT("/tags/:id", tagHandler.Update)
			admin.DELETE("/tags/:id", tagHandler.Delete)

			// Works management
			admin.POST("/works", workHandler.Create)
			admin.PUT("/works/:id", workHandler.Update)
			admin.DELETE("/works/:id", workHandler.Delete)

			// Comments management
			admin.PUT("/comments/:id/status", commentHandler.UpdateStatus)
		}
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r
}

