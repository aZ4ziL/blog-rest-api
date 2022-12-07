package models

import (
	"time"

	"gorm.io/gorm"
)

// Article
// this is the model for the article
type Article struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CategoryID  uint           `json:"category_id"`
	UserID      uint           `json:"user_id"`
	Title       string         `gorm:"size:100;unique;index" json:"title"`
	Slug        string         `gorm:"size:100;unique;index" json:"slug"`
	Logo        string         `gorm:"size:255;null" json:"logo"`
	Description string         `gorm:"size:255" json:"description"`
	Content     string         `gorm:"type:longtext" json:"content"`
	Views       uint           `json:"views"`
	Status      string         `gorm:"default:DRAFTED" json:"status"`
	Likes       []*User        `gorm:"many2many:user_article_likes" json:"likes"`
	Approved    bool           `gorm:"default:0" json:"approved"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime:nano" json:"updated_at"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Comments    []Comment      `gorm:"foreignKey:ArticleID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"comments"`
}
