package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/aZ4ziL/blog-rest-api/auth"
	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.Request.Header.Get("Authorization")
		if !strings.Contains(authorizationHeader, "Bearer") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Please login using authentication token.",
			})
			return
		}

		tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)

		claims, err := auth.ReadAndVerifyToken(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		newContext := context.WithValue(context.Background(), "user", claims)

		r := ctx.Request.WithContext(newContext)

		ctx.Request = r
		ctx.Next()
	}
}
