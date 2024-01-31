package repository

import "github.com/gin-gonic/gin"

func Index(ctx *gin.Context) {

}
func Create(ctx *gin.Context) {
	ret := map[string]interface{}{
		"Data": "ok",
	}
	ctx.JSON(200, ret)
}