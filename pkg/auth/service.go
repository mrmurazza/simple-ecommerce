package auth

import (
	"crypto/md5"
	"crypto/sha256"
	"ecommerce/config"
	"ecommerce/domain/user"
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type service struct {
	JwtSigningMethod *jwt.SigningMethodHMAC
}

func InitAuthService() *service {
	return &service{
		JwtSigningMethod: jwt.SigningMethodHS256,
	}
}

func (s *service) Authenticate(c *gin.Context) {
	cfg := config.Get()
	authorizationHeader := c.GetHeader("Authorization")
	if !strings.Contains(authorizationHeader, "Bearer") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
		return
	}

	tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method invalid")
		} else if method != s.JwtSigningMethod {
			return nil, fmt.Errorf("signing method invalid")
		}

		key := []byte(cfg.JWTSignatureKey)
		return key, nil
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.Set("userInfo", claims["data"])
}

func (s *service) AuthenticateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		s.Authenticate(c)
		c.Next()
	}
}

func (s *service) AuthenticateAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		s.Authenticate(c)
		val, exist := c.Get("userInfo")
		if !exist {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		userInfo, ok := val.(map[string]string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		if userInfo["role"] != string(user.RoleAdmin) {
			c.JSON(http.StatusForbidden, gin.H{"error": "role insufficient"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (s *service) GeneratePassword(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}

	return string(b)
}

func (s *service) EncryptPassword(password string) string {
	md5pass := md5.Sum([]byte(password))
	sha256pass := sha256.Sum256(md5pass[:])

	str := base64.StdEncoding.EncodeToString(sha256pass[:])
	return str
}

func (s *service) GetUserInfo(c *gin.Context) (map[string]interface{}, error) {
	userInfo, exist := c.Get("userInfo")
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userInfo from auth invalid"})
		c.Abort()
		return nil, errors.New("userInfo from auth invalid")
	}

	u, ok := userInfo.(map[string]interface{})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userInfo from auth invalid"})
		c.Abort()
		return nil, errors.New("userInfo from auth invalid")
	}

	return u, nil
}

func (s *service) TokenizeData(data interface{}) (string, error) {
	cfg := config.Get()
	claims := MyClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    cfg.ApplicationName,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.JWTExpiryDuration)),
		},
		Data: data,
	}

	token := jwt.NewWithClaims(
		s.JwtSigningMethod,
		claims,
	)

	key := []byte(cfg.JWTSignatureKey)
	signedToken, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
