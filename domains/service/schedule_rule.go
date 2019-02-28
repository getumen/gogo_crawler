package service

import (
	"github.com/getumen/gogo_crawler/domains/models"
	"time"
)

type ScheduleRule interface {
	UpdateStatsSuccess(request *models.Request)
	UpdateStatsFail(request *models.Request)
	ScheduleNextSuccess(request *models.Request) time.Time
	ScheduleNextFail(request *models.Request) time.Time
}
