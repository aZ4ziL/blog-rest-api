package routers

import (
	"github.com/aZ4ziL/blog-rest-api/handlers"
	"github.com/gin-gonic/gin"
)

func UserRouterV1NotAuth(group *gin.RouterGroup) {
	userHandler := handlers.NewUserHandlerV1()
	group.POST("/sign-up", userHandler.SignUp())
	group.POST("/get-token", userHandler.GetToken())
}

func UserRouterV1WithAuth(group *gin.RouterGroup) {
	userHandler := handlers.NewUserHandlerV1()

	group.GET("/user", userHandler.Auth())
}
