package logic

import (
	"EvoBot/backend/app/dto/request"
	"EvoBot/backend/app/dto/response"
	"EvoBot/backend/constant"
	"EvoBot/backend/global"
	"EvoBot/backend/utils/jwt"
	"EvoBot/backend/utils/passwd"
	"context"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type AuthLogic struct{}

type IAuthLogic interface {
	Login(c *gin.Context, req request.Login) (*response.Tokens, error)
	Logout(c *gin.Context, token string) error
}

func NewIAuthLogic() IAuthLogic {
	return &AuthLogic{}
}

func (u *AuthLogic) Login(c *gin.Context, req request.Login) (*response.Tokens, error) {
	user, err := authRepo.Get(authRepo.WithByUsername(req.Username))
	if err != nil {
		return nil, errors.WithMessage(constant.ErrRecordNotFound, err.Error())
	}
	if !passwd.Verify(req.Password, user.Password) {
		return nil, constant.ErrAuth
	}
	accessToken, err := jwt.AccessToken(user.UUID)
	if err != nil {
		global.ZAPLOG.Error("failed to generate access token", zap.Error(err))
		return nil, err
	}
	redisKey := constant.KeyUserTokenPrefix + user.UUID + ":" + c.RemoteIP()
	expireDuration := global.CONF.AuthConfig.JwtExpired
	if expireDuration <= 0 {
		expireDuration = 2 * time.Hour
	}
	err = global.RDS.Set(context.Background(), redisKey, accessToken, expireDuration).Err()
	if err != nil {
		global.ZAPLOG.Error("failed to store token in Redis", zap.Error(err))
		return nil, err
	}
	return &response.Tokens{
		AccessToken: accessToken,
	}, nil
}

func (u *AuthLogic) Logout(c *gin.Context, token string) error {
	parts := strings.SplitN(token, " ", 2)
	claims, err := jwt.VerifyToken(parts[1])
	if err != nil {
		global.ZAPLOG.Error("failed to parse token", zap.Error(err))
		return err
	}
	redisKey := constant.KeyUserTokenPrefix + claims.UUID + ":" + c.RemoteIP()
	err = global.RDS.Del(context.Background(), redisKey).Err()
	if err != nil {
		global.ZAPLOG.Error("failed to delete token from Redis", zap.Error(err))
		return err
	}
	return nil
}
