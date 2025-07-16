package router

import (
	v1 "EvoBot/backend/app/api/v1"
	"EvoBot/backend/middleware"

	"github.com/gin-gonic/gin"
)

type PolicyRouter struct {
}

func (s *PolicyRouter) InitRouter(Router *gin.RouterGroup) {
	policyRouter := Router.Group("policy").Use(middleware.AuthRequired())
	baseApi := v1.ApiGroupApp.BaseApi
	{
		policyRouter.GET("", baseApi.PolicyList)
		policyRouter.POST("", baseApi.PolicyAdd)
		policyRouter.GET("/:uuid", baseApi.PolicyGet)
		policyRouter.PATCH("/:uuid", baseApi.PolicyUpdate)
		policyRouter.DELETE("/:uuid", baseApi.PolicyDel)
	}
}
