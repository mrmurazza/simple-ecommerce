package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var chars = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type MyClaims struct {
	jwt.RegisteredClaims
	Data interface{} `json:"data"`
}

type Service interface {
	AuthenticateMiddleware() gin.HandlerFunc
	GeneratePassword(length int) string
	EncryptPassword(password string) string
	GetUserInfo(c *gin.Context) (map[string]interface{}, error)
	TokenizeData(data interface{}) (string, error)
}
