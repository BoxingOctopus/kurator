package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"media_collection/config"
	"media_collection/handlers"
	"media_collection/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// MediaItemRequest represents the request body for creating/updating a media item
type MediaItemRequest struct {
	Title         string    `json:"title" binding:"required"`
	Type          MediaType `json:"type" binding:"required"`
	PublisherID   uint      `json:"publisher_id"`
	ReleaseDate   string    `json:"release_date"`
	ISBNUPC       string    `json:"isbn_upc"`
	Format        string    `json:"format"`
	Condition     string    `json:"condition"`
	PurchaseDate  string    `json:"purchase_date"`
	PurchasePrice float64   `json:"purchase_price"`
	Notes         string    `json:"notes"`
	CreatorIDs    []uint    `json:"creator_ids"`
	GenreIDs      []uint    `json:"genre_ids"`
	PlatformIDs   []uint    `json:"platform_ids"`
}

// setupRoutes configures the API routes
func setupRoutes(router *gin.Engine, cfg *config.Config) {
	// Public routes
	router.POST("/api/register", handleRegister)
	router.POST("/api/login", handleLogin)

	// Protected routes
	protected := router.Group("/api")
	protected.Use(AuthMiddleware(cfg))
	{
		// Media routes
		protected.GET("/media", listMediaItems)
		protected.POST("/media", createMediaItem)
		protected.GET("/media/:id", getMediaItem)
		protected.PUT("/media/:id", updateMediaItem)
		protected.DELETE("/media/:id", deleteMediaItem)

		// Admin routes
		admin := protected.Group("/admin")
		admin.Use(AdminMiddleware())
		{
			admin.GET("/users", handleListUsers)
			admin.DELETE("/users/:id", handleDeleteUser)
		}
	}
}

// createMediaItem handles POST /api/media
func createMediaItem(c *gin.Context) {
	var req models.MediaItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")
	media := &models.MediaItem{
		Title:       req.Title,
		Type:        req.Type,
		ReleaseDate: req.ReleaseDate,
		Description: req.Description,
		Rating:      req.Rating,
		UserID:      userID,
	}

	if err := db.Create(media).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create media item"})
		return
	}

	c.JSON(http.StatusCreated, media.ToResponse())
}

// listMediaItems handles GET /api/media
func listMediaItems(c *gin.Context) {
	var media []models.MediaItem
	userID := c.GetUint("user_id")

	if err := db.Where("user_id = ?", userID).Find(&media).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list media items"})
		return
	}

	responses := make([]models.MediaItemResponse, len(media))
	for i, item := range media {
		responses[i] = item.ToResponse()
	}

	c.JSON(http.StatusOK, responses)
}

// getMediaItem handles GET /api/media/:id
func getMediaItem(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetUint("user_id")

	var media models.MediaItem
	if err := db.Where("media_id = ? AND user_id = ?", id, userID).First(&media).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "media item not found"})
		return
	}

	c.JSON(http.StatusOK, media.ToResponse())
}

// updateMediaItem handles PUT /api/media/:id
func updateMediaItem(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetUint("user_id")

	var media models.MediaItem
	if err := db.Where("media_id = ? AND user_id = ?", id, userID).First(&media).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "media item not found"})
		return
	}

	var req models.MediaItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	media.Title = req.Title
	media.Type = req.Type
	media.ReleaseDate = req.ReleaseDate
	media.Description = req.Description
	media.Rating = req.Rating

	if err := db.Save(&media).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update media item"})
		return
	}

	c.JSON(http.StatusOK, media.ToResponse())
}

// deleteMediaItem handles DELETE /api/media/:id
func deleteMediaItem(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetUint("user_id")

	result := db.Where("media_id = ? AND user_id = ?", id, userID).Delete(&models.MediaItem{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete media item"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "media item not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

// StartServer initializes and starts the API server
func StartServer() error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %v", err)
	}

	// Initialize database
	dsn := cfg.GetDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// Set database connection in handlers
	handlers.SetDB(db)

	// Auto-migrate models
	if err := db.AutoMigrate(&models.User{}, &models.Media{}); err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	router := gin.Default()
	setupRoutes(router, cfg)

	// Start server
	addr := fmt.Sprintf(":%d", cfg.APIPort)
	log.Printf("Starting server on %s", addr)
	return router.Run(addr)
}

func handleRegister(c *gin.Context) {
	var req models.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := CreateUser(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user.ToResponse())
}

func handleLogin(c *gin.Context) {
	var req models.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := GetUserByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if !CheckPassword(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	cfg, err := config.Load()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load configuration"})
		return
	}

	token, err := GenerateToken(user, cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user.ToResponse(),
	})
}

func handleListUsers(c *gin.Context) {
	limit := 10
	offset := 0

	if limitStr := c.Query("limit"); limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil && val > 0 {
			limit = val
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if val, err := strconv.Atoi(offsetStr); err == nil && val >= 0 {
			offset = val
		}
	}

	users, err := ListUsers(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list users"})
		return
	}

	userResponses := make([]models.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = user.ToResponse()
	}

	c.JSON(http.StatusOK, userResponses)
}

func handleDeleteUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	if err := DeleteUser(uint(userID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}

	c.Status(http.StatusNoContent)
}
