package models

import (
    "github.com/jinzhu/gorm"
)

// AppContext is used to provide stateful application configuration data to stateless endpoint handlers
type AppContext struct {
    // App config
    Config Config
    // Gorm Database
    Db *gorm.DB
}

type AppContextProvider interface {
    Ctx () *AppContext
}