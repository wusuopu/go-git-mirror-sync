package routes

import (
	"app/controllers/repository"

	"github.com/gin-gonic/gin"
)

func InitRepository(r *gin.RouterGroup) {
	r.GET("/", repository.Index)
	r.POST("/", repository.Create)
	r.PUT("/:id", repository.Update)
	r.DELETE("/:id", repository.Delete)
}
