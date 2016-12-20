package models

import (
    "github.com/jinzhu/gorm"
)

type AppContext struct {
    Db *gorm.DB
}