package service

import (
	"context"
	"github.com/getumen/gogo_crawler/domains/models"
	"time"
)

type ScheduleService interface {
	GenerateRequest(ctx context.Context, namespace string, out chan<- *models.Request)
	ScheduleRequest(ctx context.Context, in <-chan *models.Response)
	ScheduleNewRequest(ctx context.Context, in <-chan *models.Request)
}

type ScheduleRule interface {
	UpdateStatsSuccess(request *models.Request)
	UpdateStatsFail(request *models.Request)
	ScheduleNextSuccess(request *models.Request) time.Time
	ScheduleNextFail(request *models.Request) time.Time
}
