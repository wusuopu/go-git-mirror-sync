package routes

import (
	"app/controllers/repository"

	"github.com/gin-gonic/gin"
)

func InitRepository(r *gin.RouterGroup) {
	r.GET("/", repository.Index)
	r.POST("/", repository.Create)
}
