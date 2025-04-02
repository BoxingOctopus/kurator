package models

import (
	"time"
)

// MediaType represents the type of media
type MediaType string

const (
	// Book represents a book
	Book MediaType = "book"
	// Movie represents a movie
	Movie MediaType = "movie"
	// TVShow represents a TV show
	TVShow MediaType = "tv_show"
	// Game represents a game
	Game MediaType = "game"
	// Music represents music
	Music MediaType = "music"
)

// Media represents a media item in the collection
type Media struct {
	MediaID     uint      `gorm:"primaryKey" json:"media_id"`
	Title       string    `gorm:"not null" json:"title"`
	Type        MediaType `gorm:"type:varchar(20);not null" json:"type"`
	Author      string    `json:"author"`
	Genre       string    `json:"genre"`
	ReleaseDate time.Time `json:"release_date"`
	Rating      float32   `json:"rating"`
	Notes       string    `json:"notes"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// MediaRequest represents the request body for creating/updating media
type MediaRequest struct {
	Title       string    `json:"title" binding:"required"`
	Type        MediaType `json:"type" binding:"required,oneof=book movie tv_show game music"`
	Author      string    `json:"author"`
	Genre       string    `json:"genre"`
	ReleaseDate time.Time `json:"release_date"`
	Rating      float32   `json:"rating" binding:"min=0,max=10"`
	Notes       string    `json:"notes"`
}

// MediaItem represents a media item in the collection
type MediaItem struct {
	MediaID     uint   `gorm:"primaryKey"`
	Title       string `gorm:"not null"`
	Type        string `gorm:"not null"` // book, movie, game, etc.
	ReleaseDate time.Time
	Description string
	Rating      float32
	CreatedAt   time.Time
	UpdatedAt   time.Time
	UserID      uint `gorm:"not null"` // Owner of the media item
}

// MediaItemRequest represents the request to create or update a media item
type MediaItemRequest struct {
	Title       string    `json:"title" binding:"required"`
	Type        string    `json:"type" binding:"required"`
	ReleaseDate time.Time `json:"release_date"`
	Description string    `json:"description"`
	Rating      float32   `json:"rating"`
}

// MediaItemResponse represents the response for a media item
type MediaItemResponse struct {
	MediaID     uint      `json:"media_id"`
	Title       string    `json:"title"`
	Type        string    `json:"type"`
	ReleaseDate time.Time `json:"release_date"`
	Description string    `json:"description"`
	Rating      float32   `json:"rating"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UserID      uint      `json:"user_id"`
}

// ToResponse converts a MediaItem to a MediaItemResponse
func (m *MediaItem) ToResponse() MediaItemResponse {
	return MediaItemResponse{
		MediaID:     m.MediaID,
		Title:       m.Title,
		Type:        m.Type,
		ReleaseDate: m.ReleaseDate,
		Description: m.Description,
		Rating:      m.Rating,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		UserID:      m.UserID,
	}
}
