package ratelimit

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CacheData struct {
	ExpiredAt *time.Time
	Count     int
}

const MAX_RATE int = 100
const EXPIRY_DURATION = time.Duration(1 * time.Minute)

type service struct {
	Data map[string]*CacheData
}

func NewRateLimitSvc() service {
	return service{
		Data: make(map[string]*CacheData),
	}
}
func (s *service) RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		expiredAt := now.Add(EXPIRY_DURATION)
		val, ok := s.Data[ip]
		if !ok || val.ExpiredAt.Before(now) {
			s.Data[ip] = &CacheData{
				Count:     1,
				ExpiredAt: &expiredAt,
			}
			c.Next()
			return
		}

		if val.Count >= MAX_RATE {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "too many request"})
			c.Abort()
			return
		}

		val.Count += 1
	}
}
