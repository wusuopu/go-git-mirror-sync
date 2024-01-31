package initialize

import (
	"app/routes"

	"github.com/gin-gonic/gin"
)

func InitRoutes(e *gin.Engine) {
	v1 := e.Group("/api/v1")
	routes.Init((v1))
}