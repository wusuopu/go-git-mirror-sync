package middlewares

import (
	"app/schemas"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ErrorHandleMiddlewareFactory() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, err any) {
		if errors.Is(err.(error), gorm.ErrRecordNotFound) {
			// 数据库没有查询到数据，返回404
			schemas.MakeErrorResponse(c, "", 404)
			c.Abort()
			return
		}

		schemas.MakeErrorResponse(c, "", 500)
		c.Abort()
	})
}