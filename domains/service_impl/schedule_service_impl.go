package service_impl

import (
	"context"
	"github.com/getumen/gogo_crawler/domains/models"
	"github.com/getumen/gogo_crawler/domains/repository"
	"github.com/getumen/gogo_crawler/domains/service"
	"gonum.org/v1/gonum/stat/distuv"
	"log"
	"math"
	"net/http"
	"time"
)

const generateRequestLimit = 10
const heartBeat = time.Second

type scheduleService struct {
	requestRepository  repository.RequestRepository
	responseRepository repository.ResponseRepository
	scheduleRule       service.ScheduleRule
}

func NewScheduleService(
	requestRepository repository.RequestRepository,
	responseRepository repository.ResponseRepository,
	scheduleRule service.ScheduleRule) service.ScheduleService {
	return &scheduleService{
		requestRepository:  requestRepository,
		responseRepository: responseRepository,
		scheduleRule:       scheduleRule}
}

func (s *scheduleService) GenerateRequest(ctx context.Context, domain string, out chan<- *models.Request) {
	ticker := time.NewTicker(heartBeat)
LOOP:
	for {
		select {
		case <-ticker.C:
			requests, err := s.requestRepository.FindAllByDomainAndBeforeTimeOrderByNextRequest(
				ctx,
				domain,
				time.Now(),
				0,
				generateRequestLimit)

			if err != nil {
				log.Println(err)
				continue
			}
			log.Printf("Scheduling %d request\n", len(requests))
			for _, request := range requests {
				request.JobStatus = models.Running
				err = s.requestRepository.Save(ctx, request)
				if err != nil {
					log.Printf("Fail to construct request: %v", request)
				}

				out <- request
			}
		case <-ctx.Done():
			break LOOP
		default:
			time.Sleep(heartBeat)
		}
	}
	log.Printf("shutdown GenerateRequest")
	close(out)
}

func (s *scheduleService) ScheduleRequest(ctx context.Context, in <-chan *models.Response) {
	for resp := range in {

		r, err := s.requestRepository.FindByUrl(ctx, resp.Request.URL.String())
		if err != nil {
			log.Println(err)
			continue
		}

		if resp.StatusCode == http.StatusOK {

			exist, err := s.responseRepository.IsExist(ctx, resp)
			if err != nil {
				log.Println(err)
				continue
			}

			if exist {
				log.Printf("Not updated url: %s\n", r.Url.String())
				r.NextRequest = time.Now().Add(time.Since(r.LastRequest) * 2)
			} else {
				log.Printf("updated url: %s\n", r.Url.String())

				s.scheduleRule.UpdateStatsSuccess(r)

				r.LastRequest = time.Now()

				// schedule next request by schedule rule
				r.NextRequest = s.scheduleRule.ScheduleNextSuccess(r)

			}
		} else {
			log.Printf("http error %d in %s", resp.StatusCode, resp.Request.URL.String())
			s.scheduleRule.UpdateStatsFail(r)
			r.NextRequest = s.scheduleRule.ScheduleNextFail(r)
		}

		r.JobStatus = models.Complete
		err = s.requestRepository.Save(ctx, r)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}

func (s *scheduleService) ScheduleNewRequest(ctx context.Context, in <-chan *models.Request) {
	for request := range in {
		exist, err := s.requestRepository.IsExist(ctx, request.Url.String())
		if err != nil {
			log.Println(err)
			continue
		}
		if !exist {
			request.NextRequest = request.LastRequest.Add(time.Duration(1) * time.Minute)
			request.Stats = map[string]float64{}
			request.Stats["AccessCount"] = 1
			request.Stats["AccessSuccess"] = 1
			request.Stats["IntervalSum"] = 0
			err = s.requestRepository.Save(ctx, request)
			if err != nil {
				log.Println(err)
				continue
			} else {
				log.Printf("Schedule new request: %s\n", request.Url.String())
			}
		}
	}
}

func gammaSampling(count, intervalSum float64) float64 {
	gamma := distuv.Gamma{Alpha: count, Beta: math.Max(intervalSum, 1)}
	return gamma.Rand()
}

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
					r.Stats["AccessSuccess"],
					r.Stats["IntervalSum"])),
			r.Stats["IntervalSum"]/r.Stats["AccessSuccess"])) * time.Second
	return time.Now().Add(du)
}

func (poissonProcessRule) ScheduleNextFail(r *models.Request) time.Time {
	// wait twice interval
	return time.Now().Add(time.Since(r.LastRequest) * 2)
}

func (poissonProcessRule) UpdateStatsFail(r *models.Request) {
	r.Stats["AccessCount"]++
}

func (poissonProcessRule) UpdateStatsSuccess(r *models.Request) {
	r.Stats["AccessCount"]++
	r.Stats["AccessSuccess"]++
	r.Stats["IntervalSum"] += time.Since(r.LastRequest).Seconds()
}
