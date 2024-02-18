package repository

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
	var data []models.Repository

	pagination := helper.Pagination{}
	pagination.Build(di.Container.DB, ctx).Find(&data)
	utils.ThrowIfError(di.Container.DB.Error)

	schemas.MakeResponse(ctx, data, &pagination)
}
func Create(ctx *gin.Context) {
	body := helper.GetJSONBody(ctx)
	Username := helper.GetJSONString(body, "Username")
	Password := helper.GetJSONString(body, "Password")
	SSHKey := helper.GetJSONString(body, "SSHKey")

	fmt.Println("create with", helper.GetJSONString(body, "Name"))
	Name := helper.GetJSONString(body, "Name")
	if di.Container.DB.Where("name = ?", Name).First(&models.Repository{}).RowsAffected > 0 {
		schemas.MakeErrorResponse(ctx, "Name重复", 400)
		return
	}
	obj := models.Repository{
		Name: Name,
		Alias: helper.GetJSONString(body, "Alias"),
		Url: helper.GetJSONString(body, "Url"),
		AuthType: helper.GetJSONString(body, "AuthType"),
		Username: &Username,
		Password: &Password,
		SSHKey: &SSHKey,
	}
	di.Container.DB.Create(&obj)
	if di.Container.DB.Error != nil {
		schemas.MakeErrorResponse(ctx, di.Container.DB.Error, 500)
		return
	}

	schemas.MakeResponse(ctx, obj, nil)
}
func Show(ctx *gin.Context) {
	entity := models.Repository{}
	results := di.Container.DB.First(&entity, ctx.Param("repositoryId"))
	utils.ThrowIfError(results.Error)

	schemas.MakeResponse(ctx, entity, nil)
}

func Delete(ctx *gin.Context) {
	// 使用 Unscoped 彻底删除
	results := di.Container.DB.Unscoped().Delete(&models.Repository{}, ctx.Param("repositoryId"))
	utils.ThrowIfError(results.Error)
	schemas.MakeResponse(ctx, "ok", nil)
}
func Update(ctx *gin.Context) {
	entity := models.Repository{}
	results := di.Container.DB.First(&entity, ctx.Param("repositoryId"))
	utils.ThrowIfError(results.Error)

	body := helper.GetJSONBody(ctx)
	Username := helper.GetJSONString(body, "Username")
	Password := helper.GetJSONString(body, "Password")
	SSHKey := helper.GetJSONString(body, "SSHKey")

	obj := models.Repository{}
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

	di.Container.DB.Model(&entity).Updates(obj)
	if di.Container.DB.Error != nil {
		schemas.MakeErrorResponse(ctx, di.Container.DB.Error, 500)
		return
	}

	schemas.MakeResponse(ctx, obj, nil)
}

func Pull(ctx *gin.Context) {
	entity := models.Repository{}
	results := di.Container.DB.First(&entity, ctx.Param("repositoryId"))
	utils.ThrowIfError(results.Error)

	go func ()  {
		if !utils.FsIsExist(di.Container.RepositoryService.GetCoreStorePath(entity)) {
			err := di.Container.RepositoryService.Clone(entity)
			var data = make(map[string]interface{})
			if err != nil {
				data["LastError"] = err.Error()
			} else {
				data["InitedAt"] = time.Now()
				data["LastError"] = nil
			}
			di.Container.DB.Model(&entity).Updates(data)
		}
		if !utils.FsIsExist(di.Container.RepositoryService.GetCoreStorePath(entity)) {
			// 仓库未初始化
			return
		}
		err := di.Container.RepositoryService.SyncOrigin(entity)
		var data = make(map[string]interface{})
		if err != nil {
			data["LastError"] = err.Error()
		} else {
			data["PulledAt"] = time.Now()
			data["LastError"] = nil
		}
		di.Container.DB.Model(&entity).Updates(data)

		di.Container.RepositoryService.BuildBranchInfo(entity)
	}()

	schemas.MakeResponse(ctx, entity, nil)
}