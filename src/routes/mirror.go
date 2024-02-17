package routes

import (
	"app/controllers/mirror"

	"github.com/gin-gonic/gin"
)

func InitMirror(r *gin.RouterGroup) {
	r.GET("/", mirror.Index)
	r.POST("/", mirror.Create)
	r.PUT("/:branchId", mirror.Update)
	r.DELETE("/:branchId", mirror.Delete)
	r.PUT("/:branchId/push", mirror.Push)
}