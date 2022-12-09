package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aZ4ziL/blog-rest-api/auth"
	"github.com/aZ4ziL/blog-rest-api/models"
	"github.com/aZ4ziL/blog-rest-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type commentHandlerV1 struct{}

func NewCommentHandlerV1() commentHandlerV1 {
	return commentHandlerV1{}
}

// Add
// handler for add new comment into the article
func (c commentHandlerV1) Add() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.Request.Context().Value("user").(auth.Claims)
		if slug, ok := ctx.GetQuery("slug"); ok {
			article, err := models.NewArticleModel().GetArticleBySlug(slug)
			if err != nil {
				ctx.JSON(http.StatusOK, gin.H{
					"status":  "error",
					"message": fmt.Sprintf("category with slug: %s is not found.", slug),
				})
				return
			}

			commentPayload := &utils.CommentPayload{}
			if err := ctx.ShouldBindWith(&commentPayload, binding.FormPost); err != nil {
				ctx.JSON(http.StatusOK, gin.H{
					"status":  "error",
					"message": fmt.Sprintf("category with slug: %s is not found.", slug),
				})
				return
			}

			// Validate the request json
			validate = validator.New()
			err = validate.Struct(commentPayload)
			if err != nil {
				if _, ok := err.(*validator.InvalidValidationError); ok {
					log.Println(err.Error())
					return
				}
				errorMessages := []string{}
				for _, err := range err.(validator.ValidationErrors) {
					errorMessages = append(errorMessages, fmt.Sprintf("error on field: %s, with error type: %s.", err.Field(), err.ActualTag()))
				}
				ctx.JSON(http.StatusBadRequest, gin.H{
					"status":  "error",
					"message": errorMessages,
				})
				return
			}

			comment := models.Comment{
				ArticleID: article.ID,
				UserID:    user.Credential.ID,
				Text:      commentPayload.Text,
			}

			err = models.NewCommentModel().CreateComment(&comment)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"status":  "error",
					"message": err.Error(),
				})
				return
			}

			ctx.JSON(http.StatusCreated, gin.H{
				"status":  "success",
				"message": fmt.Sprintf("Successfully to comment artile slug: %s", article.Slug),
			})
			return
		}
	}
}

// Edit
// handler for edit a comment into the article
func (c commentHandlerV1) Edit() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.Request.Context().Value("user").(auth.Claims)

		commentPayloadEdit := &utils.CommentPayloadEdit{}
		err := ctx.ShouldBind(commentPayloadEdit)
		if err != nil {
			http.Error(ctx.Writer, err.Error(), http.StatusBadRequest)
			return
		}

		comment, err := models.NewCommentModel().GetCommentByID(commentPayloadEdit.CommentID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": fmt.Sprintf("comment with id: %d is not found.", commentPayloadEdit.CommentID),
			})
			return
		}

		// check user is the comment user
		if user.Credential.ID != comment.UserID {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "you are not the commenter of this comment.",
			})
			return
		}

		// check the text from request POST
		text := ctx.PostForm("text")

		if text != "" {
			comment.Text = text
		}

		err = models.GetDB().Save(&comment).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Successfully to update comment.",
			"comment": comment,
		})
	}
}

// Delete
// handler for delete comment
func (c commentHandlerV1) Delete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.Request.Context().Value("user").(auth.Claims)

		commentID := struct {
			CommentID uint `json:"comment_id" form:"comment_id"`
		}{}

		if err := ctx.ShouldBind(&commentID); err != nil {
			http.Error(ctx.Writer, err.Error(), http.StatusBadRequest)
			return
		}

		comment, err := models.NewCommentModel().GetCommentByID(commentID.CommentID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": fmt.Sprintf("The comment with ID: %d is not found error.", commentID.CommentID),
			})
			return
		}

		// checks if the logged in user is the comment owner user
		if user.Credential.ID != comment.UserID {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Sorry you are not allowed to delete this comment, because you are not the owner of this comment.",
			})
			return
		}

		if err = models.GetDB().Delete(&comment).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": fmt.Sprintf("The comment with ID: %d failed to delete error.", commentID.CommentID),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": fmt.Sprintf("The comment with ID: %d was successfully deleted.", commentID.CommentID),
		})
	}
}
