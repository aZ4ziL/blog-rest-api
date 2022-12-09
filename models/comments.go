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
	Likes       []*User        `gorm:"many2many:user_comment_likes" json:"likes"`
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
	Likes     []*User        `gorm:"many2many:user_sub_comment_likes" json:"likes"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime:nano" json:"updated_at"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type commentModel struct {
	db *gorm.DB
}

func NewCommentModel() *commentModel {
	return &commentModel{db}
}

// CreateComment
// creating a new comment, please enter `user_id` and `article_id`.
func (c *commentModel) CreateComment(comment *Comment) error {
	return c.db.Create(comment).Error
}

// GetCommentByArticleID
// returns all comment by `article_id`
func (c *commentModel) GetCommentByArticleID(articleID uint) []Comment {
	var comments []Comment
	c.db.Model(&Comment{}).Where("article_id = ?", articleID).
		Preload("SubComments").
		Preload("Likes").
		Find(&comments)
	return comments
}

// GetCommentByID
func (c *commentModel) GetCommentByID(id uint) (Comment, error) {
	var comment Comment
	err := c.db.Model(&Comment{}).Where("id = ?", id).Preload("SubComments").Preload("Likes").
		First(&comment).Error
	return comment, err
}

// CreateSubComment
// create new sub comment
func (c *commentModel) CreateSubComment(subComment *SubComment) error {
	return c.db.Create(subComment).Error
}
