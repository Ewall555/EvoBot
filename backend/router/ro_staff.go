package router

import (
	v1 "EvoBot/backend/app/api/v1"
	"EvoBot/backend/middleware"

	"github.com/gin-gonic/gin"
)

type StaffRouter struct {
}

func (s *StaffRouter) InitRouter(Router *gin.RouterGroup) {
	staffRouter := Router.Group("staff").Use(middleware.AuthRequired())
	baseApi := v1.ApiGroupApp.BaseApi
	{
		staffRouter.GET("", baseApi.StaffList)
		staffRouter.POST("", baseApi.StaffAdd)
		staffRouter.GET("/:uuid", baseApi.StaffGet)
		staffRouter.PATCH("/:uuid", baseApi.StaffUpdate)
		staffRouter.DELETE("/:uuid", baseApi.StaffDel)
	}
}
