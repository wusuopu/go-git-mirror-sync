package initialize

import (
	"app/config"

	"github.com/gin-gonic/gin"
)

func commonInit(e *gin.Engine) *gin.Engine {
	config.Load()
	InitServices()

	var engine *gin.Engine
	if e == nil {
		engine = gin.Default()
	} else {
		engine = e
	}
	InitDB()
	InitRoutes(engine)
	return engine
}

func Init(e *gin.Engine) *gin.Engine {
	// 先加载 .env 文件
	InitEnv()

	return commonInit(e)
}

func InitTest(e *gin.Engine) *gin.Engine {
	// 先加载 .env.test 文件
	InitEnv(".env.test")
	config.Config["GO_ENV"] = "test"

	gin.SetMode(gin.TestMode)
	return commonInit(e)
}
