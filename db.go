package main

import (
	"fmt"
	"time"

	"media_collection/config"
	"media_collection/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// MediaType represents the type of media
type MediaType string

const (
	VideoGame MediaType = "video_game"
	Music     MediaType = "music"
	Book      MediaType = "book"
	Movie     MediaType = "movie"
)

// Publisher represents a publisher/studio/label
type Publisher struct {
	PublisherID uint   `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Website     string
	CreatedAt   time.Time
}

// Creator represents an author, artist, director, or developer
type Creator struct {
	CreatorID uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	Type      string
	CreatedAt time.Time
}

// MediaItem represents a media item (book, movie, game, etc.)
type MediaItem struct {
	ItemID        uint      `gorm:"primaryKey"`
	Title         string    `gorm:"not null"`
	Type          MediaType `gorm:"type:media_type;not null"`
	PublisherID   uint
	Publisher     Publisher
	ReleaseDate   time.Time
	ISBNUPC       string
	Format        string
	Condition     string
	PurchaseDate  time.Time
	PurchasePrice float64
	Notes         string
	CreatedAt     time.Time
	Creators      []Creator  `gorm:"many2many:media_creators;"`
	Genres        []Genre    `gorm:"many2many:media_genres;"`
	Platforms     []Platform `gorm:"many2many:game_platforms;"`
}

// Genre represents a media genre
type Genre struct {
	GenreID   uint      `gorm:"primaryKey"`
	Name      string    `gorm:"not null"`
	MediaType MediaType `gorm:"type:media_type;not null"`
	CreatedAt time.Time
}

// Platform represents a gaming platform
type Platform struct {
	PlatformID uint   `gorm:"primaryKey"`
	Name       string `gorm:"not null"`
}

// MediaCreator represents the junction table for media items and creators
type MediaCreator struct {
	MediaID   uint
	CreatorID uint
	Role      string
	MediaItem MediaItem
	Creator   Creator
}

// MediaGenre represents the junction table for media items and genres
type MediaGenre struct {
	MediaID   uint
	GenreID   uint
	MediaItem MediaItem
	Genre     Genre
}

// GamePlatform represents the junction table for games and platforms
type GamePlatform struct {
	MediaID    uint
	PlatformID uint
	MediaItem  MediaItem
	Platform   Platform
}

// User represents a user in the system
type User struct {
	UserID    uint   `gorm:"primaryKey"`
	Username  string `gorm:"not null"`
	Email     string `gorm:"not null"`
	CreatedAt time.Time
}

var db *gorm.DB

// InitDB initializes the database connection and creates tables
func InitDB(cfg *config.Config) error {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(&models.MediaItem{}, &models.User{})
	if err != nil {
		return fmt.Errorf("failed to migrate database schema: %w", err)
	}

	return nil
}

// CRUD operations for MediaItem

// CreateMediaItem creates a new media item
func CreateMediaItem(item *MediaItem) error {
	return db.Create(item).Error
}

// GetMediaItem retrieves a media item by ID
func GetMediaItem(id uint) (*MediaItem, error) {
	var item MediaItem
	err := db.Preload("Publisher").
		Preload("Creators").
		Preload("Genres").
		Preload("Platforms").
		First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// UpdateMediaItem updates an existing media item
func UpdateMediaItem(item *MediaItem) error {
	return db.Save(item).Error
}

// DeleteMediaItem deletes a media item
func DeleteMediaItem(id uint) error {
	return db.Delete(&MediaItem{}, id).Error
}

// ListMediaItems retrieves all media items with optional filters
func ListMediaItems(mediaType *MediaType, limit, offset int) ([]MediaItem, error) {
	var items []MediaItem
	query := db.Preload("Publisher").
		Preload("Creators").
		Preload("Genres").
		Preload("Platforms")

	if mediaType != nil {
		query = query.Where("type = ?", *mediaType)
	}

	err := query.Limit(limit).Offset(offset).Find(&items).Error
	return items, err
}

// Example usage:
// func main() {
// 	// Initialize database connection
// 	err := InitDB(&config.Config{
// 		DBHost:     "localhost",
// 		DBUser:     "postgres",
// 		DBPassword: "your_password",
// 		DBName:     "media_collection",
// 		DBPort:     5432,
// 	})
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Example: Create a new media item
// 	newItem := &MediaItem{
// 		Title:    "The Great Gatsby",
// 		Type:     Book,
// 		ISBNUPC:  "978-0743273565",
// 		Format:   "Paperback",
// 		Condition: "Good",
// 	}

// 	err = CreateMediaItem(newItem)
// 	if err != nil {
// 		log.Printf("Error creating media item: %v", err)
// 	}

// 	// Example: List all books
// 	bookType := Book
// 	books, err := ListMediaItems(&bookType, 10, 0)
// 	if err != nil {
// 		log.Printf("Error listing books: %v", err)
// 	}

// 	for _, book := range books {
// 		fmt.Printf("Book: %s\n", book.Title)
// 	}
// }
