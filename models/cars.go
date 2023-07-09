package models

import "gorm.io/gorm"

type Cars struct {
	ID   uint    `gorm:"primary key;autoIncrement" json:"id"`
	Name *string `json:"name"`
	Year *string `json:"year"`
	Type *string `json:"type"`
}

func MigrateCars(db *gorm.DB) error {
	err := db.AutoMigrate(&Cars{})
	return err
}
