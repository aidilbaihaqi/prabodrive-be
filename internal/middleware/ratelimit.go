package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"

	"github.com/aidilbaihaqi/prabodrive-be/internal/delivery/http/response"
)

func RateLimit() gin.HandlerFunc {
	rate := limiter.Rate{
		Period: time.Minute,
		Limit:  100,
	}

	store := memory.NewStore()
	instance := limiter.New(store, rate)

	middleware := mgin.NewMiddleware(instance, mgin.WithLimitReachedHandler(func(c *gin.Context) {
		response.TooManyRequests(c)
		c.Abort()
	}))

	return middleware
}
