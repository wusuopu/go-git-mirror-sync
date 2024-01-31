package repository

import (
	"app/controllers/repository"

	"github.com/gin-gonic/gin"
)

func Init(r *gin.RouterGroup) {
	r.GET("/", repository.Index)
	r.POST("/", repository.Create)
}