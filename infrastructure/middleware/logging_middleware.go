package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("call to endpoints: %s %s", c.Request.Method, c.Request.URL.Path)
		c.Next()
	}
}
