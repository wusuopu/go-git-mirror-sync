package initialize

import (
	"app/config"

	"github.com/gin-gonic/gin"
)

func Init(e *gin.Engine) *gin.Engine {
	// 先加载 .env 文件
	InitEnv()

	var engine *gin.Engine
	if e == nil {
		engine = gin.Default()
	} else {
		engine = e
	}
	config.Load()
	InitRoutes(engine)
	return engine
}