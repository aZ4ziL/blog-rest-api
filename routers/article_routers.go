package routers

import (
	"github.com/aZ4ziL/blog-rest-api/handlers"
	"github.com/gin-gonic/gin"
)

func ArticleRouterV1(group *gin.RouterGroup) {
	articleHandler := handlers.NewArticleHandlerV1()
	commentHandler := handlers.NewCommentHandlerV1()

	group.Any("", articleHandler.All())
	group.POST("/comment/add", commentHandler.Add())
	group.PUT("/comment/edit", commentHandler.Edit())
	group.DELETE("/comment/delete", commentHandler.Delete())
}
