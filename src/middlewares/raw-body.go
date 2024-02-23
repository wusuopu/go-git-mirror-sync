package middlewares

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
)

func RawBodyMiddleware(c *gin.Context) {
	if (c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH") {
		if strings.HasPrefix(c.Request.Header.Get("Content-Type"), "application/json") {
			rawBody, err := c.GetRawData()
			if err == nil {
				// 获取原始的 body 内容
				c.Set("rawBody", rawBody)
				c.Request.Body = io.NopCloser(bytes.NewBuffer(rawBody))
			} else {
				fmt.Println(err)
			}
		}
	}
	c.Next()
}