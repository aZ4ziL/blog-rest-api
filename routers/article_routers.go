package routers

import (
	"github.com/aZ4ziL/blog-rest-api/handlers"
	"github.com/gin-gonic/gin"
)

func ArticleRouterV1(group *gin.RouterGroup) {
	articleHandler := handlers.NewArticleHandlerV1()
	group.Any("", articleHandler.All())
}
