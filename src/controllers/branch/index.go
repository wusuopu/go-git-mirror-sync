package branch

import (
	"app/di"
	"app/models"
	"app/schemas"

	"github.com/gin-gonic/gin"
)

func Index(ctx *gin.Context) {
	var data []models.Branch

	di.Container.DB.Where("repository_id = ?", ctx.Param("repositoryId")).Find(&data)
	schemas.MakeResponse(ctx, data, nil)
}