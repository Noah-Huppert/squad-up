package db

type User struct {
    TableMetadata
    ID int `gorm:"serial primary key"`
    FirstName string
    LastName string
    Email string
    ProfilePictureUrl string
}

