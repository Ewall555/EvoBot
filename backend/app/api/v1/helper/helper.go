package helper

import (
	"EvoBot/backend/app/dto"
	"EvoBot/backend/constant"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SuccessWithData(ctx *gin.Context, data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	res := dto.Response{
		Code: constant.CodeSuccess,
		Data: data,
	}
	ctx.JSON(http.StatusOK, res)
	ctx.Abort()
}

func SuccessWithOutData(ctx *gin.Context) {
	res := dto.Response{
		Code:    constant.CodeSuccess,
		Message: "success",
	}
	ctx.JSON(http.StatusOK, res)
	ctx.Abort()
}

func SuccessWithMsg(ctx *gin.Context, msg string) {
	res := dto.Response{
		Code:    constant.CodeSuccess,
		Message: msg,
	}
	ctx.JSON(http.StatusOK, res)
	ctx.Abort()
}

func ErrResponse(ctx *gin.Context, code int) {
	ctx.JSON(code, nil)
	ctx.Abort()
}

func ErrResponseWithErr(ctx *gin.Context, code int, err error) {
	res := dto.Response{
		Code:    code,
		Message: err.Error(),
	}
	ctx.JSON(code, res)
	ctx.Abort()
}

func ErrResponseWithMsg(ctx *gin.Context, code int, msg string) {
	res := dto.Response{
		Code:    code,
		Message: msg,
	}
	ctx.JSON(code, res)
	ctx.Abort()
}
