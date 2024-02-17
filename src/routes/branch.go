package routes

import (
	"app/controllers/branch"

	"github.com/gin-gonic/gin"
)

func InitBranch(r *gin.RouterGroup) {
	r.GET("/", branch.Index)
}
