package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aZ4ziL/blog-rest-api/auth"
	"github.com/aZ4ziL/blog-rest-api/models"
	"github.com/aZ4ziL/blog-rest-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// userHandlerV1
// user handler V1
type userHandlerV1 struct{}

func NewUserHandlerV1() userHandlerV1 {
	return userHandlerV1{}
}

/*
SignUp
Handler for user registration

This handler uses POST requests. Where the payload uses JSON Payload.

Example:

	{
		"first_name": "First Name",
		"last_name": "Last Name",
		"username": "Username",
		"email": "Email Address",
		"password1": "Password"
		"password2": "Password confirmation"
	}
*/
func (u userHandlerV1) SignUp() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userPayload := &utils.UserSignUpPayload{}

		// decode the body to json payload
		err := ctx.ShouldBindJSON(userPayload)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Failed to decode request to json",
			})
			return
		}

		// Validate the request json
		validate = validator.New()
		err = validate.Struct(userPayload)
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

		// check confirmation password
		if userPayload.Password1 != userPayload.Password2 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Please enter the same confirmation password.",
			})
			return
		}

		// Save the user
		user := models.User{
			FirstName: userPayload.FirstName,
			LastName:  userPayload.LastName,
			Username:  userPayload.Username,
			Email:     userPayload.Email,
			Password:  userPayload.Password2,
		}
		err = models.NewUserModel().CreateUser(&user)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		// If no error return response 201
		ctx.JSON(http.StatusCreated, gin.H{
			"status":  "success",
			"message": "Successfully to create new user by Username: " + user.Username,
		})
	}
}

/*
GetToken
handler to get a new token.
This handler uses the post method where the payload is of type JSON.

	{
		"username": "your username",
		"password": "your password"
	}

If the username and password are valid then you will get a token.
*/
func (u userHandlerV1) GetToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userPayload := &utils.UserSignInPayload{}
		// decode request to json payload
		err := ctx.ShouldBindJSON(userPayload)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Failed to decode request to json",
			})
			return
		}

		// Validate the request json
		validate = validator.New()
		err = validate.Struct(userPayload)
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

		// Check the username
		user, err := models.NewUserModel().GetUserByUsername(userPayload.Username)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Username or password is incorrect.",
			})
			return
		}
		// Check password
		if !auth.DecryptionPassword(user.Password, userPayload.Password) {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Username or password is incorrect.",
			})
			return
		}

		// manual converting to Credential
		credential := auth.Credential{
			ID:          user.ID,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Username:    user.Username,
			Email:       user.Email,
			IsSuperuser: user.IsSuperuser,
			IsStaff:     user.IsStaff,
			IsActive:    user.IsActive,
			LastLogin:   user.LastLogin.Time,
			DateJoined:  user.DateJoined,
		}

		token, err := auth.GenerateNewToken(credential)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Failed to generate new token, with error message: " + err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Please do not share your tokens with anyone.",
			"token":   token,
		})
	}
}

// Auth
func (u userHandlerV1) Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get user
		user := ctx.Request.Context().Value("user").(auth.Claims)
		ctx.JSON(http.StatusOK, user)
	}
}
