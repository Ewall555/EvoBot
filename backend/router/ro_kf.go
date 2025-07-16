package router

import (
	v1 "EvoBot/backend/app/api/v1"
	"EvoBot/backend/middleware"

	"github.com/gin-gonic/gin"
)

type KFRouter struct {
}

func (s *KFRouter) InitRouter(Router *gin.RouterGroup) {
	kfRouter := Router.Group("kf").Use(middleware.AuthRequired())
	baseApi := v1.ApiGroupApp.BaseApi
	{
		kfRouter.GET("", baseApi.KFList)
		kfRouter.POST("", baseApi.KFAdd)
		kfRouter.GET("/:uuid", baseApi.KFGet)
		kfRouter.PATCH("/:uuid", baseApi.KFUpdate)
		kfRouter.DELETE("/:uuid", baseApi.KFDel)
	}
}
