package models

import (
	"database/sql"
	"time"

	"github.com/aZ4ziL/blog-rest-api/auth"
	"gorm.io/gorm"
)

// Declare a model for the user.
type User struct {
	ID              uint          `gorm:"primaryKey" json:"id"`
	FirstName       string        `gorm:"size:30" json:"first_name"`
	LastName        string        `gorm:"size:30" json:"last_name"`
	Username        string        `gorm:"size:30;unique;index" json:"username"`
	Email           string        `gorm:"size:30;unique;index" json:"email"`
	Password        string        `gorm:"size:100" json:"-"`
	IsSuperuser     bool          `gorm:"default:0" json:"is_superuser"`
	IsStaff         bool          `gorm:"default:0" json:"is_staff"`
	IsActive        bool          `gorm:"default:1" json:"is_active"`
	LastLogin       sql.NullTime  `json:"last_login"`
	DateJoined      time.Time     `gorm:"autoCreateTime" json:"date_joined"`
	Articles        []Article     `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"articles"`
	Comments        []Comment     `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"comments"`
	SubComments     []SubComment  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"sub_comments"`
	ArticleLikes    []*Article    `gorm:"many2many:user_article_likes" json:"articles_likes"`
	CommentLikes    []*Comment    `gorm:"many2many:user_comment_likes" json:"comment_likes"`
	SubCommentLikes []*SubComment `gorm:"many2many:user_sub_comment_likes" json:"sub_comment_likes"`
}

// userModel
type userModel struct {
	db *gorm.DB
}

func NewUserModel() *userModel {
	return &userModel{db}
}

// CreateUser
// create a new user, if the user already exists it
// will return an error message `user with user/email already in use`.
func (u *userModel) CreateUser(user *User) error {
	user.Password = auth.EncryptionPassword(user.Password) // encrypt the password
	err := u.db.Create(user).Error
	return err
}

// GetUserByID
// this function will return a user object by passing the `id` parameter.
// If the user with the requested `id` is not found it will return an error message.
func (u *userModel) GetUserByID(id uint) (User, error) {
	var user User
	err := u.db.Model(&User{}).Where("id = ?", id).
		Preload("Articles").
		Preload("Comments").
		Preload("SubComments").
		Preload("ArticleLikes").
		Preload("CommentLikes").
		Preload("SubCommentLikes").
		First(&user).Error
	return user, err
}

// GetUserByUsername
// this function will return a user object by passing the `username` parameter.
// If the user with the requested `username` is not found it will return an error message.
func (u *userModel) GetUserByUsername(username string) (User, error) {
	var user User
	err := u.db.Model(&User{}).Where("username = ?", username).
		Preload("Articles").
		Preload("Comments").
		Preload("SubComments").
		Preload("ArticleLikes").
		Preload("CommentLikes").
		Preload("SubCommentLikes").
		First(&user).Error
	return user, err
}

// GetUserByEmail
// this function will return a user object by passing the `email` parameter.
// If the user with the requested `email` is not found it will return an error message.
func (u *userModel) GetUserByEmail(email string) (User, error) {
	var user User
	err := u.db.Model(&User{}).Where("email = ?", email).
		Preload("Articles").
		Preload("Comments").
		Preload("SubComments").
		Preload("ArticleLikes").
		Preload("CommentLikes").
		Preload("SubCommentLikes").
		First(&user).Error
	return user, err
}
