package routes

import (
	"github.com/gin-gonic/gin"
)

func Init(router *gin.RouterGroup, engine *gin.Engine) {
	InitRepository(router.Group("/repositories"))
	engine.GET("_health", func(ctx *gin.Context) {
		ctx.String(200, "ok")
	})
}
