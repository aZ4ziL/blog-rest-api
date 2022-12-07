package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/aZ4ziL/blog-rest-api/auth"
	"github.com/aZ4ziL/blog-rest-api/models"
	"github.com/aZ4ziL/blog-rest-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type categoryHandler struct{}

func NewCategoryHandler() categoryHandler {
	return categoryHandler{}
}

// Index
// get all categories
func (c categoryHandler) Index() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if slug, ok := ctx.GetQuery("slug"); ok {
			category, err := models.NewCategoryModel().GetCategoryBySlug(slug)
			if err != nil {
				ctx.JSON(http.StatusOK, gin.H{
					"status":  "error",
					"message": fmt.Sprintf("category with slug: %s is not found.", slug),
				})
				return
			}

			ctx.JSON(http.StatusOK, category)
			return
		} else {
			categories := models.NewCategoryModel().GetAllCategories()
			ctx.JSON(http.StatusOK, categories)
			return
		}
	}
}

// Add
func (c categoryHandler) Add() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// This is only the super and staff user
		user := ctx.Request.Context().Value("user").(auth.Claims)
		if !user.IsSuperuser && !user.IsStaff {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "You don't have permission to access this method.",
			})
			return
		}

		categoryPayload := &utils.CategoryPayload{}

		err := ctx.ShouldBindWith(categoryPayload, binding.FormMultipart)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		// Validate the request json
		validate = validator.New()
		err = validate.Struct(categoryPayload)
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

		file, err := ctx.FormFile("logo")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
		}

		contentType := file.Header.Get("Content-Type")
		log.Println(contentType)
		if contentType != "image/jpg" && contentType != "image/png" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Please upload image with format `png|jpg`",
			})
			return
		}

		filename := "media/categories/" + uuid.NewString() + file.Filename

		category := models.Category{
			Title:    categoryPayload.Title,
			Slug:     categoryPayload.Slug,
			Logo:     filename,
			Approved: categoryPayload.Approved,
		}

		err = models.NewCategoryModel().CreateCategory(&category)
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
			"message": "Successfully to add new category.",
			"url":     "http://" + ctx.Request.Host + "/v1/categories?slug=" + category.Slug,
		})
	}
}

// Update
func (c categoryHandler) Update() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// This is only the super and staff user
		user := ctx.Request.Context().Value("user").(auth.Claims)
		if !user.IsSuperuser && !user.IsStaff {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "You don't have permission to access this method.",
			})
			return
		}

		if slug, ok := ctx.GetQuery("slug"); ok {
			category, err := models.NewCategoryModel().GetCategoryBySlug(slug)
			if err != nil {
				ctx.JSON(http.StatusOK, gin.H{
					"status":  "error",
					"message": fmt.Sprintf("category with slug: %s is not found.", slug),
				})
				return
			}

			title := ctx.PostForm("title")
			slug := ctx.PostForm("slug")
			approved := ctx.PostForm("approved")
			if title != "" {
				category.Title = title
			}
			if slug != "" {
				category.Slug = slug
			}
			if approved != "" {
				approvedBool, _ := strconv.ParseBool(approved)
				category.Approved = approvedBool
			}

			file, err := ctx.FormFile("logo")
			if err != nil {
				if err = models.GetDB().Save(&category).Error; err != nil {
					ctx.JSON(http.StatusOK, gin.H{
						"status":  "error",
						"message": err.Error(),
					})
					return
				} else {
					ctx.JSON(http.StatusOK, gin.H{
						"status":   "success",
						"message":  "Success updating the category with slug: " + slug,
						"category": category,
					})
					return
				}
			}

			// remove old file
			err = os.Remove(category.Logo)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, nil)
				return
			}

			filename := "media/categories/" + uuid.NewString() + file.Filename

			category.Logo = filename

			if err = models.GetDB().Save(&category).Error; err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"status":  "error",
					"message": err.Error(),
				})
				return
			}

			ctx.SaveUploadedFile(file, filename)
			ctx.JSON(http.StatusOK, gin.H{
				"status":   "success",
				"message":  "Success updating the category with slug: " + slug,
				"category": category,
			})
			return
		}
	}
}

// Delete
func (c categoryHandler) Delete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// This is only the super and staff user
		user := ctx.Request.Context().Value("user").(auth.Claims)
		if !user.IsSuperuser && !user.IsStaff {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "You don't have permission to access this method.",
			})
			return
		}

		if slug, ok := ctx.GetQuery("slug"); ok {
			category, err := models.NewCategoryModel().GetCategoryBySlug(slug)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"status":  "error",
					"message": "You don't have permission to access this method.",
				})
				return
			}

			err = models.GetDB().Unscoped().Delete(&category).Error
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"status":  "error",
					"message": "Failed deleting category with slug: " + slug,
				})
				return
			}

			// Delete file
			err = os.Remove(category.Logo)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, nil)
				return
			}

			ctx.JSON(http.StatusOK, gin.H{
				"status":  "warning",
				"message": "Successfully to delete category with slug: " + slug,
			})
			return
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "If you want to delete the category, please set the query url with key `slug`.",
			})
			return
		}
	}
}
