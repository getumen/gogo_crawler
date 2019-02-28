package service_impl

import (
	"github.com/getumen/gogo_crawler/domains/models"
	"github.com/getumen/gogo_crawler/domains/service"
	"gonum.org/v1/gonum/stat/distuv"
	"math"
	"time"
)

type poissonProcessRule struct {
}

func NewPoissonProcessRule() service.ScheduleRule {
	return poissonProcessRule{}
}

func (poissonProcessRule) ScheduleNextSuccess(r *models.Request) time.Time {
	du := time.Duration(
		math.Max(
			math.Ceil(
				gammaSampling(
					r.Stats("AccessSuccess"),
					r.Stats("IntervalSum"))),
			r.Stats("IntervalSum")/r.Stats("AccessSuccess"))) * time.Second
	return time.Now().Add(du)
}

func (poissonProcessRule) ScheduleNextFail(r *models.Request) time.Time {
	// wait twice interval
	return time.Now().Add(time.Since(r.LastRequest()) * 2)
}

func (poissonProcessRule) UpdateStatsFail(r *models.Request) {
	r.SetStats("AccessCount", r.Stats("AccessCount")+1)
}

func (poissonProcessRule) UpdateStatsSuccess(r *models.Request) {
	r.SetStats("AccessCount", r.Stats("AccessCount")+1)
	r.SetStats("AccessSuccess", r.Stats("AccessSuccess")+1)
	r.SetStats("IntervalSum", r.Stats("IntervalSum")+time.Since(r.LastRequest()).Seconds())
}

func gammaSampling(count, intervalSum float64) float64 {
	gamma := distuv.Gamma{Alpha: count, Beta: math.Max(intervalSum, 1)}
	return gamma.Rand()
}
