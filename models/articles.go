package models

import (
	"errors"
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

// declate the article models
type articleModel struct {
	db *gorm.DB
}

func NewArticleModel() *articleModel {
	return &articleModel{db}
}

// CreateArticle
// create a new article, if the article exists it will return an error message.
func (a *articleModel) CreateArticle(article *Article) error {
	return a.db.Create(article).Error
}

// GetArticleByID
// returns articles by `id`, if article `id` is not found it will return an error.
func (a *articleModel) GetArticleByID(id uint) (Article, error) {
	var article Article
	err := a.db.Model(&Article{}).Where("id = ?", id).
		Preload("Likes").Preload("Comments").First(&article).Error
	return article, err
}

// GetArticleBySlug
// returns articles by `slug`, if article `slug` is not found it will return an error.
func (a *articleModel) GetArticleBySlug(slug string) (Article, error) {
	var article Article
	err := a.db.Model(&Article{}).Where("slug = ?", slug).
		Preload("Likes").Preload("Comments").First(&article).Error
	return article, err
}

// GetAllArticles
// get all articles
func (a *articleModel) GetAllArticles() []Article {
	var articles []Article
	a.db.Model(&Article{}).
		Preload("Likes").
		Preload("Comments").
		Find(&articles)
	return articles
}

// AddLikeToArticle
// add like to article
func (a *articleModel) AddLikeToArticle(userID, articleID uint) error {
	userModel := NewUserModel()

	user, err := userModel.GetUserByID(userID)
	if err != nil {
		return err
	}

	article, err := a.GetArticleByID(articleID)
	if err != nil {
		return err
	}

	var isJoined bool
	for _, userLike := range article.Likes {
		if userLike.ID == user.ID {
			isJoined = true
		}
	}

	if !isJoined {
		a.db.Model(&article).Association("Likes").Append([]*User{&user})
	} else {
		return errors.New("you has already liked this article")
	}

	return nil
}

// DeleteLikeToArticle
func (a *articleModel) DeleteLikeToArticle(userID, articleID uint) error {
	userModel := NewUserModel()

	user, err := userModel.GetUserByID(userID)
	if err != nil {
		return err
	}

	article, err := a.GetArticleByID(articleID)
	if err != nil {
		return err
	}

	var isJoined bool
	for _, userLike := range article.Likes {
		if userLike.ID == user.ID {
			isJoined = true
		}
	}

	if isJoined {
		a.db.Model(&article).Association("Likes").Delete([]*User{&user})
	} else {
		return errors.New("you has not already liked this article")
	}

	return nil
}
