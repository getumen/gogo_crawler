package service

import (
	"context"
	"github.com/getumen/gogo_crawler/domains/models"
)

type ItemService interface {
	SaveResponse(ctx context.Context, in <-chan *models.Response)
}
