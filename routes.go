package main

import (
	"media_collection/handlers"
	"media_collection/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures the API routes
func SetupRoutes(r *gin.Engine) {
	// Public routes
	r.POST("/api/register", handlers.Register)
	r.POST("/api/login", handlers.Login)

	// Protected routes
	auth := r.Group("/api")
	auth.Use(middleware.AuthRequired())
	{
		// User routes
		auth.GET("/user/profile", handlers.GetProfile)
		auth.PUT("/user/profile", handlers.UpdateProfile)
		auth.POST("/user/password", handlers.ChangePassword)

		// 2FA routes
		auth.POST("/user/2fa/setup", handlers.SetupTOTP)
		auth.POST("/user/2fa/enable", handlers.EnableTOTP)
		auth.POST("/user/2fa/disable", handlers.DisableTOTP)

		// Media routes
		auth.GET("/media", handlers.ListMedia)
		auth.POST("/media", handlers.CreateMedia)
		auth.GET("/media/:id", handlers.GetMedia)
		auth.PUT("/media/:id", handlers.UpdateMedia)
		auth.DELETE("/media/:id", handlers.DeleteMedia)
	}
}
