package models

import (
	"gorm.io/gorm"
)

type File struct {
	ID       uint   `gorm:"primaryKey"`
	Filename string `gorm:"not null"`
	MimeType string `gorm:"not null"`
	Filepath string `gorm:"not null"`
}

type ResponseUploadFile struct {
	Url string `json:"url"`
}

func MigrateFiles(db *gorm.DB) error {
	return db.AutoMigrate(&File{})
}
