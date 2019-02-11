package repository

import (
	"context"
	"github.com/getumen/gogo_crawler/domains/models"
	"time"
)

type RequestRepository interface {
	IsExist(ctx context.Context, url string) (bool, error)
	FindAllByDomainAndBeforeTimeOrderByNextRequest(ctx context.Context, domain string, now time.Time, offset, limit int) ([]*models.Request, error)
	FindByUrl(ctx context.Context, url string) (*models.Request, error)
	Save(ctx context.Context, r *models.Request) error
}
