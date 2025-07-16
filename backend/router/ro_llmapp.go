package router

import (
	v1 "EvoBot/backend/app/api/v1"
	"EvoBot/backend/middleware"

	"github.com/gin-gonic/gin"
)

type LLMAppRouter struct {
}

func (s *LLMAppRouter) InitRouter(Router *gin.RouterGroup) {
	llmappRouter := Router.Group("llmapp").Use(middleware.AuthRequired())
	baseApi := v1.ApiGroupApp.BaseApi
	{
		llmappRouter.POST("/config", baseApi.LLMAppConfigAdd)
		llmappRouter.GET("/config", baseApi.LLMAppConfigList)
		llmappRouter.GET("/config/:uuid", baseApi.LLMAppConfigGet)
		llmappRouter.PATCH("/config/:uuid", baseApi.LLMAppConfigUpdate)
		llmappRouter.DELETE("/config/:uuid", baseApi.LLMAppConfigDel)
	}
}
