package mirror

import (
	"app/di"
	"app/models"
	"app/schemas"
	"app/utils"
	"app/utils/helper"
	"fmt"
	"time"

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
	Name := helper.GetJSONString(body, "Name")
	if di.Container.DB.Where("name = ?", Name).Where("repository_id = ?", ctx.Param("repositoryId")).First(&models.Mirror{}).RowsAffected > 0 {
		schemas.MakeErrorResponse(ctx, "Name重复", 400)
		return
	}

	Username := helper.GetJSONString(body, "Username")
	Password := helper.GetJSONString(body, "Password")
	SSHKey := helper.GetJSONString(body, "SSHKey")

	obj := models.Mirror{
		Name: Name,
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
	err := di.Container.RepositoryService.CreateRemote(entity, obj)
	if err != nil {
		di.Container.Logger.Error(fmt.Sprintf("Create remote mirror error: %s %s %s", entity.Name, obj.Name, err.Error()))
	}

	schemas.MakeResponse(ctx, obj, nil)
}
func Update(ctx *gin.Context) {
	entity := models.Repository{}
	results := di.Container.DB.First(&entity, ctx.Param("repositoryId"))
	utils.ThrowIfError(results.Error)

	oldObj := models.Mirror{}
	results = di.Container.DB.First(&oldObj, ctx.Param("mirrorId"))
	utils.ThrowIfError(results.Error)

	body := helper.GetJSONBody(ctx)
	Username := helper.GetJSONString(body, "Username")
	Password := helper.GetJSONString(body, "Password")
	SSHKey := helper.GetJSONString(body, "SSHKey")

	needUpdate := false

	obj := models.Mirror{}
	Alias := helper.GetJSONString(body, "Alias")
	Url := helper.GetJSONString(body, "Url")
	AuthType := helper.GetJSONString(body, "AuthType")
	if oldObj.Alias != Alias {
		obj.Alias = Alias
	}
	if oldObj.Url != Url {
		// remote 的 URL 有变化需要更新仓库配置
		needUpdate = true
		obj.Url = Url
		oldObj.Url = Url
	}
	if oldObj.AuthType != AuthType {
		obj.AuthType = AuthType
	}
	if *oldObj.Username != Username {
		obj.Username = &Username
	}
	if *oldObj.Password != Password {
		obj.Password = &Password
	}
	if *oldObj.SSHKey != SSHKey {
		obj.SSHKey = &SSHKey
	}
	results = di.Container.DB.Model(&oldObj).Updates(obj)
	utils.ThrowIfError(results.Error)


	if needUpdate {
		// 更新仓库的 remote
		err := di.Container.RepositoryService.CreateRemote(entity, oldObj)
		if err != nil {
			di.Container.Logger.Error(fmt.Sprintf("Create remote mirror error: %s %s %s", entity.Name, obj.Name, err.Error()))
		}
	}

	schemas.MakeResponse(ctx, obj, nil)
}
func Delete(ctx *gin.Context) {
	entity := models.Repository{}
	results := di.Container.DB.First(&entity, ctx.Param("repositoryId"))
	utils.ThrowIfError(results.Error)

	obj := models.Mirror{}
	results = di.Container.DB.First(&obj, ctx.Param("mirrorId"))
	utils.ThrowIfError(results.Error)

	// 使用 Unscoped 彻底删除
	results = di.Container.DB.Unscoped().Delete(&models.Mirror{}, ctx.Param("mirrorId"))
	utils.ThrowIfError(results.Error)

	// 在仓库创建 remote
	di.Container.RepositoryService.DeleteRemote(entity, obj)

	schemas.MakeResponse(ctx, "ok", nil)
}
func Push(ctx *gin.Context) {
	entity := models.Mirror{}
	results := di.Container.DB.Preload("Repository").First(&entity, ctx.Param("mirrorId"))
	utils.ThrowIfError(results.Error)

	go func ()  {
		err := di.Container.RepositoryService.SyncMirror(entity.Repository, entity)
		var data = make(map[string]interface{})
		if err != nil {
			data["LastError"] = err.Error()
		} else {
			data["PushedAt"] = time.Now()
			data["LastError"] = nil
		}
		di.Container.DB.Model(&entity).Updates(data)
	}()
	schemas.MakeResponse(ctx, "ok", nil)
}