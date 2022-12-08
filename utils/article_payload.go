package utils

type ArticlePayload struct {
	CategoryID  uint   `form:"category_id" validate:"required,number"`
	UserID      uint   `form:"user_id" validate:"required,number"`
	Title       string `form:"title" validate:"required"`
	Slug        string `form:"slug" validate:"max=30"`
	Description string `form:"description" validate:"required,max=255"`
	Content     string `form:"content" validate:"required"`
}
