package utils

type CommentPayload struct {
	Text string `form:"text" validate:"required"`
}

type CommentPayloadEdit struct {
	CommentID uint `form:"comment_id" validate:"required,number"`
}
