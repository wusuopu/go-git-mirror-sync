package initialize

import (
	"app/middlewares"
	"app/routes"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/itchyny/timefmt-go"
)

func InitRoutes(e *gin.Engine) {
	e.Use(gin.LoggerWithFormatter(func (param gin.LogFormatterParams) string {
		headers := "{"
		for k, v := range param.Request.Header {
			line := ""
			for _, item := range v {
				line = line + item + ";"
			}
			headers = headers + k + ":" + line + " "
		}
		headers = headers + "}"

		return fmt.Sprintf("%s - [%s] %s %s %s %d %s \"%s\" \"%s\"\n",
				param.ClientIP,
				timefmt.Format(param.TimeStamp, "%Y-%m-%d %H:%M:%S %z"),
				param.Method,
				param.Path,
				param.Request.Proto,
				param.StatusCode,
				param.Latency,
				headers,
				param.ErrorMessage,
		)
	}))
	e.Use(gin.Recovery())

	v1 := e.Group("/api/v1", middlewares.BasicAuthMiddleware, middlewares.RawBodyMiddleware)
	routes.Init(v1, e)
}
