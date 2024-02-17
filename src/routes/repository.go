package routes

import (
	"app/controllers/repository"

	"github.com/gin-gonic/gin"
)

func InitRepository(r *gin.RouterGroup) {
	r.GET("/", repository.Index)
	r.POST("/", repository.Create)
	r.GET("/:repositoryId", repository.Show)
	r.PUT("/:repositoryId", repository.Update)
	r.DELETE("/:repositoryId", repository.Delete)
	r.PUT("/:repositoryId/pull", repository.Pull)
}
