package models

import "gorm.io/gorm"

type User struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `gorm:"not null" json:"password"`
	Role     string `gorm:"not null" json:"role"`
}

func MigrateUser(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}