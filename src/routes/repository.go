package routes

import (
	"app/controllers/repository"

	"github.com/gin-gonic/gin"
)

func InitRepository(r *gin.RouterGroup) {
	r.GET("/", repository.Index)
	r.POST("/", repository.Create)
	r.GET("/:id", repository.Show)
	r.PUT("/:id", repository.Update)
	r.DELETE("/:id", repository.Delete)
	r.PUT("/:id/pull", repository.Pull)
}
