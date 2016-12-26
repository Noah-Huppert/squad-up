package db

import "time"

// Basic lifecycle metadata fields provided by gorm.Model but without the ID field. This allows for customizations of
// the ID field such as autoincrement.
type TableMetadata struct {
    // Time row was created at
	CreatedAt time.Time `json:"created_at"`
    // Last time row was updated
	UpdatedAt time.Time `json:"updated_at"`
    // If using soft delete, time the item was deleted
	DeletedAt *time.Time `json:"deleted_at"`
}

