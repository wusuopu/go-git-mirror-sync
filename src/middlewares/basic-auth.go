package middlewares

import (
	"os"

	"github.com/gin-gonic/gin"
)

var basicAuth gin.HandlerFunc

func BasicAuthMiddleware(c *gin.Context) {
	user := os.Getenv("BASIC_AUTH_USER")
	password := os.Getenv("BASIC_AUTH_PASSWORD")
	if user == "" || password == "" {
		c.Next()
		return
	}
	if basicAuth == nil {
		basicAuth = gin.BasicAuth(gin.Accounts{
			os.Getenv("BASIC_AUTH_USER"): os.Getenv("BASIC_AUTH_PASSWORD"),
		})
	}
	basicAuth(c)
}