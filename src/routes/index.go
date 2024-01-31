package routes

import (
	"app/routes/repository"

	"github.com/gin-gonic/gin"
)

func Init(router *gin.RouterGroup) {
	repository.Init(router.Group("/repositories"))
}