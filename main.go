package main

import (
	"github.com/aZ4ziL/blog-rest-api/middlewares"
	"github.com/aZ4ziL/blog-rest-api/routers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"})

	// user group not auth
	userGroupV1NoAuth := r.Group("/v1/auth")
	routers.UserRouterV1NotAuth(userGroupV1NoAuth)

	// user group with auth
	userGroupV1WithAuth := r.Group("/v1/auth")
	userGroupV1WithAuth.Use(middlewares.Authentication())
	routers.UserRouterV1WithAuth(userGroupV1WithAuth)

	r.Run(":8000")
}
