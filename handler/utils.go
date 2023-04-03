package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type Context *gin.Context

func fetchUserId(c *gin.Context) (*int, error) {
	userCtx, exist := c.Get("userInfo")
	if !exist {
		return nil, errors.New("unauthorized")
	}
	userInfo, ok := userCtx.(map[string]interface{})
	if !ok {
		return nil, errors.New("malformed auth metadata")
	}

	userId := int((userInfo["id"]).(float64))
	return &userId, nil
}
