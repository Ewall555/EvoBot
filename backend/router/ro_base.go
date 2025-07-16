package router

import (
	v1 "EvoBot/backend/app/api/v1"

	"github.com/gin-gonic/gin"
)

type BaseRouter struct{}

func (s *BaseRouter) InitRouter(Router *gin.RouterGroup) {
	baseRouter := Router.Group("auth")
	baseApi := v1.ApiGroupApp.BaseApi
	{
		baseRouter.POST("/login", baseApi.Login)
		baseRouter.POST("/logout", baseApi.Logout)
	}
}
