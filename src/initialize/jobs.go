package initialize

import (
	"app/config"
	"app/di"
	"app/jobs"
	"fmt"

	"github.com/go-co-op/gocron/v2"
)

func InitJobs() {
	s, err := gocron.NewScheduler()
	if err == nil {
		di.Container.Scheduler = &s
		fmt.Println("Run CronJob in ", config.Config["CRONTAB"].(string))
		s.NewJob(
			gocron.CronJob(config.Config["CRONTAB"].(string), false),
			gocron.NewTask(jobs.RepositorySyncJob),
		)
	}
}