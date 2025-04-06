package middleware

import (
	"time"

	"github.com/axiaoxin-com/ratelimiter"
	"github.com/gin-gonic/gin"
)

func RateLimitMiddleware() gin.HandlerFunc {
	return ratelimiter.GinMemRatelimiter(ratelimiter.GinRatelimiterConfig{
		LimitKey: func(c *gin.Context) string {
			return c.ClientIP()
		},
		LimitedHandler: func(c *gin.Context) {
			c.JSON(200, "too many requests!!!")
			c.Abort()
			return
		},
		TokenBucketConfig: func(*gin.Context) (time.Duration, int) {
			return time.Second * 60, 4000
		},
	})
}
