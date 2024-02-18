package initialize

import (
	"app/config"
	"app/di"
	"app/jobs"

	"github.com/go-co-op/gocron/v2"
)

func InitJobs() {
	s, err := gocron.NewScheduler()
	if err == nil {
		di.Container.Scheduler = &s
		s.NewJob(
			gocron.CronJob(config.Config["CRONTAB"].(string), false),
			gocron.NewTask(jobs.RepositorySyncJob),
		)
	}
}