package router

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (r Router) Handler() *gin.Engine {
	engine := gin.Default()

	// Public group (no auth)
	public := engine.Group("/api/public")
	{
		public.GET("/ping", func(c *gin.Context) { c.JSON(200, gin.H{"msg": "pong"}) })
	}

	// Authenticated group (JWT required)
	auth := engine.Group("/api")
	auth.Use(jwtAuthMiddleware())
	{
		auth.GET("/audit-logs", r.authenticatedRESTHandler.GetAuditLogs)
	}

	return engine
}

// jwtAuthMiddleware is a simple middleware that checks for Bearer token in Authorization header
func jwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") || len(authHeader) <= 7 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			return
		}
		token := strings.TrimSpace(authHeader[7:])
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing JWT token"})
			return
		}
		// Optionally: validate JWT here
		c.Next()
	}
}
