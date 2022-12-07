package models

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	ArticleID   uint           `json:"article_id"`
	UserID      uint           `json:"user_id"`
	Text        string         `gorm:"type:longtext" json:"text"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime:nano" json:"updated_at"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	SubComments []SubComment   `gorm:"foreignKey:CommentID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"sub_comments"`
}

type SubComment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CommentID uint           `json:"comment_id"`
	UserID    uint           `json:"user_id"`
	Text      string         `gorm:"type:longtext" json:"text"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime:nano" json:"updated_at"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
