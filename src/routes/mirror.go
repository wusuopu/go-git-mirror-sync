package routes

import (
	"app/controllers/mirror"

	"github.com/gin-gonic/gin"
)

func InitMirror(r *gin.RouterGroup) {
	r.GET("/", mirror.Index)
	r.POST("/", mirror.Create)
	r.PUT("/:mirrorId", mirror.Update)
	r.DELETE("/:mirrorId", mirror.Delete)
	r.PUT("/:mirrorId/push", mirror.Push)
}