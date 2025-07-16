package jwt

import (
	"EvoBot/backend/constant"
	"EvoBot/backend/global"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type CustomClaims struct {
	UUID string `json:"uuid"`
	jwt.RegisteredClaims
}

var secret string
var secretOnce sync.Once

func initSecret() {
	authConfig := global.CONF.AuthConfig
	if authConfig == nil {
		global.ZAPLOG.Error("auth config is nil")
		return
	}
	secret = authConfig.Secret
}

func getSecret() []byte {
	secretOnce.Do(initSecret)
	return []byte(secret)
}

func generateToken(uuid string, expire time.Duration) (string, error) {
	mySigningKey := getSecret()
	claims := CustomClaims{
		UUID: uuid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "evobot",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func AccessToken(uuid string) (string, error) {
	expireDuration := global.CONF.AuthConfig.JwtExpired
	if expireDuration <= 0 {
		expireDuration = 2 * time.Hour
	}
	return generateToken(uuid, expireDuration)
}

func VerifyToken(tokenString string) (*CustomClaims, error) {
	mySigningKey := getSecret()
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})
	if err != nil || token == nil {
		return nil, constant.ErrTokenParse
	}
	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, constant.ErrTokenParse
	}
	return claims, nil
}
