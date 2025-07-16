package router

import (
	v1 "EvoBot/backend/app/api/v1"
	"EvoBot/backend/middleware"

	"github.com/gin-gonic/gin"
)

type WecomRouter struct {
}

func (s *WecomRouter) InitRouter(Router *gin.RouterGroup) {
	wecomRouter := Router.Group("wecom")
	wecomAuthRouter := Router.Group("wecom").Use(middleware.AuthRequired())
	baseApi := v1.ApiGroupApp.BaseApi
	{
		wecomRouter.POST("/callback", baseApi.WecomHandle)
		wecomRouter.GET("/callback", baseApi.WecomHandle)

		wecomAuthRouter.GET("/config", baseApi.WecomConfigList)
		wecomAuthRouter.GET("/config/:uuid", baseApi.WecomConfigGet)
		wecomAuthRouter.PATCH("/config/:uuid", baseApi.WecomConfigUpdate)

		wecomAuthRouter.GET("/receptionist/:kfid", baseApi.WecomReceptionistList)
		wecomAuthRouter.POST("/receptionist/:kfid", baseApi.WecomReceptionistAdd)
		wecomAuthRouter.DELETE("/receptionist/:kfid", baseApi.WecomReceptionistDel)

		wecomAuthRouter.GET("/account", baseApi.WecomAccountList)
		wecomAuthRouter.GET("/account/:kfid", baseApi.WecomAddContactWay)
	}
}
