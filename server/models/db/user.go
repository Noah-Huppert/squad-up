package db

type User struct {
    TableMetadata
    ID int `gorm:"serial primary key" json:"id"`
    FirstName string  `json:"first_name"`
    LastName string `json:"last_name"`
    Email string `json:"email"`
    ProfilePictureUrl string `json:"profile_picture_url"`
}

