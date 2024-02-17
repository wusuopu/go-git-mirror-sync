package mirror

import (
	"app/di"
	"app/models"
	"app/schemas"
	"app/utils"
	"app/utils/helper"
	"fmt"

	"github.com/gin-gonic/gin"
)

func Index(ctx *gin.Context) {
	var data []models.Mirror

	di.Container.DB.Where("repository_id = ?", ctx.Param("repositoryId")).Find(&data)
	schemas.MakeResponse(ctx, data, nil)
}

func Create(ctx *gin.Context) {
	entity := models.Repository{}
	results := di.Container.DB.First(&entity, ctx.Param("repositoryId"))
	utils.ThrowIfError(results.Error)

	body := helper.GetJSONBody(ctx)
	Username := helper.GetJSONString(body, "Username")
	Password := helper.GetJSONString(body, "Password")
	SSHKey := helper.GetJSONString(body, "SSHKey")

	obj := models.Mirror{
		Name: helper.GetJSONString(body, "Name"),
		Alias: helper.GetJSONString(body, "Alias"),
		Url: helper.GetJSONString(body, "Url"),
		AuthType: helper.GetJSONString(body, "AuthType"),
		Username: &Username,
		Password: &Password,
		SSHKey: &SSHKey,
		RepositoryId: entity.ID,
	}

	results = di.Container.DB.Create(&obj)
	utils.ThrowIfError(results.Error)

	// 在仓库创建 remote

	schemas.MakeResponse(ctx, obj, nil)
}
func Update(ctx *gin.Context) {
	entity := models.Mirror{}
	results := di.Container.DB.First(&entity, ctx.Param("mirrorId"))
	utils.ThrowIfError(results.Error)

	body := helper.GetJSONBody(ctx)
	Username := helper.GetJSONString(body, "Username")
	Password := helper.GetJSONString(body, "Password")
	SSHKey := helper.GetJSONString(body, "SSHKey")

	obj := models.Mirror{}
	Alias := helper.GetJSONString(body, "Alias")
	Url := helper.GetJSONString(body, "Url")
	AuthType := helper.GetJSONString(body, "AuthType")
	if entity.Alias != Alias {
		obj.Alias = Alias
	}
	if entity.Url != Url {
		obj.Url = Url
	}
	if entity.AuthType != AuthType {
		obj.AuthType = AuthType
	}
	if *entity.Username != Username {
		obj.Username = &Username
	}
	if *entity.Password != Password {
		obj.Password = &Password
	}
	if *entity.SSHKey != SSHKey {
		obj.SSHKey = &SSHKey
	}

	results = di.Container.DB.Model(&entity).Updates(obj)
	utils.ThrowIfError(results.Error)

	// 更新仓库的 remote

	schemas.MakeResponse(ctx, obj, nil)
}
func Delete(ctx *gin.Context) {
	// 使用 Unscoped 彻底删除
	results := di.Container.DB.Unscoped().Delete(&models.Mirror{}, ctx.Param("mirrorId"))
	utils.ThrowIfError(results.Error)
	schemas.MakeResponse(ctx, "ok", nil)
}
func Push(ctx *gin.Context) {
	entity := models.Mirror{}
	results := di.Container.DB.Preload("Repository").First(&entity, ctx.Param("mirrorId"))
	utils.ThrowIfError(results.Error)

	fmt.Println(entity)
	schemas.MakeResponse(ctx, entity, nil)
}