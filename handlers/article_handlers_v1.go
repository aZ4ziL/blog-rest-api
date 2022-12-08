package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/aZ4ziL/blog-rest-api/auth"
	"github.com/aZ4ziL/blog-rest-api/models"
	"github.com/aZ4ziL/blog-rest-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type articleHandlerV1 struct{}

func NewArticleHandlerV1() articleHandlerV1 {
	return articleHandlerV1{}
}

// Index
func (a articleHandlerV1) All() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// get the user from authentication header context
		user := ctx.Request.Context().Value("user").(auth.Claims)
		// if method is post
		// this is handler for return all articles and article detail
		if ctx.Request.Method == "GET" {
			// Query Slug
			if slug, ok := ctx.GetQuery("slug"); ok {
				article, err := models.NewArticleModel().GetArticleBySlug(slug)
				if err != nil {
					ctx.JSON(http.StatusNotFound, gin.H{
						"status":  "error",
						"message": fmt.Sprintf("the article by slug: %s is not found.", slug),
					})
					return
				}

				// add views to 1
				views := article.Views + 1
				models.GetDB().Model(&article).Update("views", views)

				ctx.JSON(http.StatusOK, article)
				return
			} else {
				articleJSON := []models.Article{}
				articles := models.NewArticleModel().GetAllArticles()

				for _, a := range articles {
					if a.Status == "PUBLISHED" && a.Approved {
						articleJSON = append(articleJSON, a)
					}
				}
				ctx.JSON(http.StatusOK, articleJSON)
				return
			}
		}

		// if method is post
		// for add new article
		if ctx.Request.Method == "POST" {
			articlePayload := &utils.ArticlePayload{}

			err := ctx.ShouldBindWith(articlePayload, binding.FormMultipart)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"status":  "error",
					"message": err.Error(),
				})
				return
			}

			// validate the form
			validate = validator.New()
			err = validate.Struct(articlePayload)
			if err != nil {
				if _, ok := err.(*validator.InvalidValidationError); ok {
					log.Println(err.Error())
					return
				}
				errorMessages := []string{}
				for _, err := range err.(validator.ValidationErrors) {
					errorMessages = append(errorMessages, fmt.Sprintf("error on field: %s, with error type: %s", err.Field(), err.ActualTag()))
				}
				ctx.JSON(http.StatusBadRequest, gin.H{
					"status":  "error",
					"message": errorMessages,
				})
				return
			}
			if articlePayload.Slug != "" {
				if strings.Contains(articlePayload.Slug, " ") {
					articlePayload.Slug = strings.ToLower(strings.ReplaceAll(articlePayload.Slug, " ", "-"))
				} else if strings.Contains(articlePayload.Slug, "-") {
					articlePayload.Slug = strings.ToLower(articlePayload.Slug)
				}
			} else {
				articlePayload.Slug = strings.ToLower(strings.ReplaceAll(articlePayload.Title, " ", "-"))
			}

			article := models.Article{
				CategoryID:  articlePayload.CategoryID,
				UserID:      articlePayload.UserID,
				Title:       articlePayload.Title,
				Slug:        articlePayload.Slug,
				Description: articlePayload.Description,
				Content:     articlePayload.Content,
			}

			file, err := ctx.FormFile("logo")
			if err != nil {
				// if user not set the logo image
				err = models.NewArticleModel().CreateArticle(&article)
				if err != nil {
					ctx.JSON(http.StatusBadRequest, gin.H{
						"status":  "error",
						"message": err.Error(),
					})
					return
				}
				ctx.JSON(http.StatusCreated, gin.H{
					"status":  "success",
					"message": "Successfully to create new article.",
				})
				return
			} else {
				// if user upload the image for logo article
				filename := "media/articles/" + uuid.NewString() + file.Filename

				article.Logo = filename
				err = models.NewArticleModel().CreateArticle(&article)
				if err != nil {
					ctx.JSON(http.StatusBadRequest, gin.H{
						"status":  "error",
						"message": err.Error(),
					})
					return
				}

				ctx.SaveUploadedFile(file, filename)

				ctx.JSON(http.StatusCreated, gin.H{
					"status":  "success",
					"message": "Successfully to create new article.",
				})
				return
			}
		}

		// if method is PUT
		// handler for edit the article by slug
		if ctx.Request.Method == "PUT" {
			// only accept the query `slug`
			if slug, ok := ctx.GetQuery("slug"); ok {
				article, err := models.NewArticleModel().GetArticleBySlug(slug)
				if err != nil {
					ctx.JSON(http.StatusNotFound, gin.H{
						"status":  "error",
						"message": fmt.Sprintf("the article by slug: %s is not found.", slug),
					})
					return
				}

				// check this user are the author or not from this article
				if user.Credential.ID != article.UserID {
					ctx.JSON(http.StatusBadRequest, gin.H{
						"status":  "permission_danied",
						"message": "error: you don't have permission to access this method.",
					})
					return
				}

				// declare and get post method value
				categoryID := ctx.PostForm("category_id")
				title := ctx.PostForm("title")
				slug := ctx.PostForm("slug")
				description := ctx.PostForm("description")
				content := ctx.PostForm("content")

				// check all method value
				if categoryID != "" {
					categoryIDInt, _ := strconv.Atoi(categoryID)
					article.CategoryID = uint(categoryIDInt)
				}
				if title != "" {
					article.Title = title
				}
				// set slug
				if slug != "" {
					if strings.Contains(slug, "-") {
						article.Slug = strings.ToLower(slug)
					}
					if strings.Contains(slug, " ") {
						article.Slug = strings.ToLower(strings.ReplaceAll(slug, " ", "-"))
					}
				}
				if description != "" {
					article.Description = description
				}
				if content != "" {
					article.Content = content
				}

				file, err := ctx.FormFile("logo")
				if err != nil {
					// if user(author) is not edit the logo of article
					err := models.GetDB().Save(&article).Error
					if err != nil {
						ctx.JSON(http.StatusBadRequest, gin.H{
							"status":  "error",
							"message": err.Error(),
						})
						return
					} else {
						ctx.JSON(http.StatusOK, gin.H{
							"status":  "success",
							"message": "Successfully to update the article with slug: " + slug,
							"article": article,
						})
						return
					}
				} else {
					_ = os.Remove(article.Logo)
					filename := "media/articles/" + uuid.NewString() + file.Filename

					article.Logo = filename
					if err := models.GetDB().Save(&article).Error; err != nil {
						ctx.JSON(http.StatusBadRequest, gin.H{
							"status":  "error",
							"message": err.Error(),
						})
						return
					} else {
						ctx.SaveUploadedFile(file, filename)
						ctx.JSON(http.StatusOK, gin.H{
							"status":  "success",
							"message": "Successfully to update the article with slug: " + slug,
							"article": article,
						})
						return
					}
				}
			}
		}

		// if method is DELETE
		// handler for delete the article by slug
		if ctx.Request.Method == "DELETE" {
			if slug, ok := ctx.GetQuery("slug"); ok {
				article, err := models.NewArticleModel().GetArticleBySlug(slug)
				if err != nil {
					ctx.JSON(http.StatusNotFound, gin.H{
						"status":  "error",
						"message": fmt.Sprintf("the article by slug: %s is not found.", slug),
					})
					return
				}

				// check user is the author
				if user.Credential.ID != article.UserID {
					ctx.JSON(http.StatusBadRequest, gin.H{
						"status":  "permission_danied",
						"message": "error: you don't have permission to access this method.",
					})
					return
				}

				// Delete it
				err = models.GetDB().Delete(&article).Error
				if err != nil {
					ctx.JSON(http.StatusBadRequest, gin.H{
						"status":  "error",
						"message": "Failed to delete article by slug: " + slug,
					})
					return
				}

				ctx.JSON(http.StatusOK, gin.H{
					"status":  "success",
					"message": "Successfully to delete article by slug: " + slug,
				})
				return
			}
		}
	}
}
