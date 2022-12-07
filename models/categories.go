package models

import (
	"time"

	"gorm.io/gorm"
)

// Category
// this type will declare which categories are available in the blog application.
type Category struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"size:30;unique;index" json:"title"`
	Slug      string         `gorm:"size:100;unique;index" json:"slug"`
	Logo      string         `gorm:"size:255;null" json:"logo"`
	Approved  bool           `gorm:"default:0" json:"approved"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime:nano" json:"updated_at"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Articles  []Article      `gorm:"foreignKey:CategoryID" json:"articles"`
}
