package routers

import (
	"github.com/aZ4ziL/blog-rest-api/handlers"
	"github.com/gin-gonic/gin"
)

func CategoryRouter(group *gin.RouterGroup) {
	categoryHandler := handlers.NewCategoryHandler()
	group.GET("", categoryHandler.Index())
	group.POST("", categoryHandler.Add())
	group.PUT("", categoryHandler.Update())
	group.DELETE("", categoryHandler.Delete())
}
