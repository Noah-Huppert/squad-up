package db

import "time"

// Basic lifecycle metadata fields provided by gorm.Model but without the ID field. This allows for customizations of
// the ID field such as autoincrement.
type TableMetadata struct {
    // Time row was created at
	CreatedAt time.Time
    // Last time row was updated
	UpdatedAt time.Time
    // If using soft delete, time the item was deleted
	DeletedAt *time.Time
}

