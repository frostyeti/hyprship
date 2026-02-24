package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/frostyeti/hyprship/apps/api/internal/routes"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type rateLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	limiters = make(map[string]*rateLimiter)
	mu       sync.Mutex
)

// cleanup removes old limiters that haven't been seen recently
func init() {
	go func() {
		for {
			time.Sleep(10 * time.Minute)
			mu.Lock()
			for ip, v := range limiters {
				if time.Since(v.lastSeen) > 30*time.Minute {
					delete(limiters, ip)
				}
			}
			mu.Unlock()
		}
	}()
}

func getLimiter(ip string, r rate.Limit, b int) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	v, exists := limiters[ip]
	if !exists {
		limiter := rate.NewLimiter(r, b)
		limiters[ip] = &rateLimiter{limiter, time.Now()}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

// RateLimit creates a middleware that restricts requests based on IP address
// reqsPerSecond is the number of allowed requests per second, and burst is the burst size.
func RateLimit(reqsPerSecond float64, burst int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getLimiter(ip, rate.Limit(reqsPerSecond), burst)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, routes.ErrorResponse(&routes.ApiError{
				Code:    "too_many_requests",
				Message: "Rate limit exceeded. Please try again later.",
			}))
			c.Abort()
			return
		}
		c.Next()
	}
}
