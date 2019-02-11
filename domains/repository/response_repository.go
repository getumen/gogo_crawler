package repository

import (
	"context"
	"github.com/getumen/gogo_crawler/domains/models"
)

type ResponseRepository interface {
	Save(ctx context.Context, response *models.Response) error
	IsExist(ctx context.Context, response *models.Response) (bool, error)
}
