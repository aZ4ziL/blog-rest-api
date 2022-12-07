package models

import (
	"database/sql"
	"time"
)

// Declare a model for the user.
type User struct {
	ID           uint         `gorm:"primaryKey" json:"id"`
	FirstName    string       `gorm:"size:30" json:"first_name"`
	LastName     string       `gorm:"size:30" json:"last_name"`
	Username     string       `gorm:"size:30;unique;index" json:"username"`
	Email        string       `gorm:"size:30;unique;index" json:"email"`
	Password     string       `gorm:"size:100" json:"-"`
	IsSuperuser  bool         `gorm:"default:0" json:"is_superuser"`
	IsStaff      bool         `gorm:"default:0" json:"is_staff"`
	IsActive     bool         `gorm:"default:1" json:"is_active"`
	LastLogin    sql.NullTime `json:"last_login"`
	DateJoined   time.Time    `gorm:"autoCreateTime" json:"date_joined"`
	Articles     []Article    `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"articles"`
	Comments     []Comment    `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"comments"`
	SubComments  []SubComment `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"sub_comments"`
	ArticleLikes []*Article   `gorm:"many2many:user_article_likes" json:"articles_likes"`
}
