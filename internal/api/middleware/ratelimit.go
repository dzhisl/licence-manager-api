package middleware

import (
	"net"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type ClientLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.Mutex
	r        rate.Limit
	b        int
}

func NewClientLimiter(r rate.Limit, b int) *ClientLimiter {
	return &ClientLimiter{
		limiters: make(map[string]*rate.Limiter),
		r:        r,
		b:        b,
	}
}

func (cl *ClientLimiter) getLimiter(ip string) *rate.Limiter {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	limiter, exists := cl.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(cl.r, cl.b)
		cl.limiters[ip] = limiter
	}

	return limiter
}

func RateLimitMiddleware(cl *ClientLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip rate limiting for /api/metrics
		if c.Request.URL.Path == "/api/metrics" {
			c.Next()
			return
		}

		ip := clientIP(c)
		limiter := cl.getLimiter(ip)

		if !limiter.Allow() {
			c.AbortWithStatusJSON(429, gin.H{
				"error": "Rate limit exceeded",
			})
			return
		}

		c.Next()
	}
}

func clientIP(c *gin.Context) string {
	ip := c.ClientIP()
	// Clean IPv6 prefix
	if ip4 := net.ParseIP(ip).To4(); ip4 != nil {
		return ip4.String()
	}
	return ip
}
