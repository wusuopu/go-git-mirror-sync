package jobs

import (
	"app/di"
	"app/models"
	"app/utils"
	"fmt"
	"time"
)

func RepositorySyncJob() {
	var repositories []models.Repository
	var mirrors []models.Mirror

	di.Container.Logger.Info("Start RepositorySyncJob")

	results := di.Container.DB.Find(&repositories)
	if results.Error != nil {
		di.Container.Logger.Error(fmt.Sprintf("RepositorySyncJob find Repository error: %s", results.Error.Error()))
		return
	}
	for _, entity := range repositories {
		// 拉取代码
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


		//推送代码
		results := di.Container.DB.Where("repository_id = ?", entity.ID).Find(&mirrors)
		if results.Error != nil {
			di.Container.Logger.Error(fmt.Sprintf("RepositorySyncJob find Mirror error for %d: %s", entity.ID, results.Error.Error()))
			continue
		}

		for _, obj := range mirrors {
			di.Container.RepositoryService.CreateRemote(entity, obj)

			err := di.Container.RepositoryService.SyncMirror(entity, obj)
			var data = make(map[string]interface{})
			if err != nil {
				data["LastError"] = err.Error()
			} else {
				data["PushedAt"] = time.Now()
				data["LastError"] = nil
			}
			di.Container.DB.Model(&obj).Updates(data)
		}
	}

	di.Container.Logger.Info("Finish RepositorySyncJob")
}