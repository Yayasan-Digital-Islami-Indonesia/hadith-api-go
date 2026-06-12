package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func RateLimit(rate string) gin.HandlerFunc {
	store := memory.NewStore()
	rateLimiter := limiter.New(store, limiter.Rate{Limit: 100, Period: 60})

	return func(c *gin.Context) {
		ctx, err := rateLimiter.Get(c, c.ClientIP())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "rate limiter error"})
			c.Abort()
			return
		}

		c.Writer.Header().Set("X-RateLimit-Limit", "100")
		c.Writer.Header().Set("X-RateLimit-Remaining", strconv.Itoa(int(ctx.Remaining)))

		if ctx.Reached {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			c.Abort()
			return
		}

		c.Next()
	}
}