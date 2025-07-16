package middleware

import (
	"EvoBot/backend/app/api/v1/helper"
	"EvoBot/backend/constant"
	"EvoBot/backend/global"
	"EvoBot/backend/utils/jwt"
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AuthRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get(constant.TokenKey)
		if token == "" {
			helper.ErrResponseWithMsg(ctx, constant.CodeErrUnauthorized, "no login")
			return
		}
		parts := strings.SplitN(token, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			helper.ErrResponseWithMsg(ctx, constant.CodeErrUnauthorized, "token format is incorrect")
			return
		}
		claims, err := jwt.VerifyToken(parts[1])
		if err != nil {
			global.ZAPLOG.With(zap.String("token", parts[1]), zap.String("IP", ctx.RemoteIP())).Error(err.Error())
			helper.ErrResponseWithMsg(ctx, constant.CodeErrUnauthorized, "invalid or expired token")
			return
		}
		if !validateTokenWithRedis(claims.UUID, parts[1], ctx.RemoteIP()) {
			helper.ErrResponseWithMsg(ctx, constant.CodeErrUnauthorized, "token expired or used elsewhere")
			return
		}
		ctx.Set("user", map[string]interface{}{
			"uuid": claims.UUID,
		})
		ctx.Next()
	}
}

func validateTokenWithRedis(uuid string, token, ip string) bool {
	result, err := global.RDS.Get(context.Background(), constant.KeyUserTokenPrefix+uuid+":"+ip).Result()
	if err != nil || result != token {
		global.ZAPLOG.With(zap.String("UUID", uuid), zap.String("IP", ip)).Error(err.Error())
		return false
	}
	return true
}
