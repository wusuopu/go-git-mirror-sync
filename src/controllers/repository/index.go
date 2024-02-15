package repository

import (
	"app/di"
	"app/models"
	"app/schemas"
	"app/utils"
	"app/utils/helper"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	obj := models.Repository{
		Name: helper.GetJSONString(body, "Name"),
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


func Delete(ctx *gin.Context) {
	// 使用 Unscoped 彻底删除
	di.Container.DB.Unscoped().Delete(&models.Repository{}, ctx.Param("id"))
	if di.Container.DB.Error != nil {
		schemas.MakeErrorResponse(ctx, di.Container.DB.Error, 500)
		return
	}
	schemas.MakeResponse(ctx, "ok", nil)
}
func Update(ctx *gin.Context) {
	entity := models.Repository{}
	di.Container.DB.Find(&entity, ctx.Param("id"))
	if errors.Is(di.Container.DB.Error, gorm.ErrRecordNotFound) {
		schemas.MakeErrorResponse(ctx, "", 404)
		return
	}

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