package main

import (
	"gorm.io/gorm"
)

type Subscription struct {
	ID          string `gorm:"primaryKey"`
	ServiceName string `json:"service_name" binding:"required"`
	Price       int    `json:"price" binding:"required"`
	UserID      string `json:"user_id" binding:"required"`
	StartDate   string `json:"start_date" binding:"required"`
	EndDate     string `json:"end_date" binding:"omitempty"`
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&Subscription{})
}
