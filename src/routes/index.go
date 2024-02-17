package routes

import (
	"app/middlewares"
	"time"

	"github.com/gin-gonic/gin"
)

func Init(router *gin.RouterGroup, engine *gin.Engine) {
	// 静态文件
	engine.Static("/statics", "./assets/statics")
	engine.LoadHTMLFiles("./assets/index.html")

	// engine.StaticFile("/", "./assets/index.html")
	engine.GET("/", middlewares.BasicAuthMiddleware, func(ctx *gin.Context) {
		ctx.HTML(200, "index.html", gin.H{
			"debug": gin.Mode() != gin.ReleaseMode,
			"version": time.Now().Unix(),
		})
	})

	InitRepository(router.Group("/repositories"))
	InitBranch(router.Group("/repositories/:repositoryId/branches"))
	InitBranch(router.Group("/repositories/:repositoryId/mirrors"))
	engine.GET("_health", func(ctx *gin.Context) {
		ctx.String(200, "ok")
	})
}
