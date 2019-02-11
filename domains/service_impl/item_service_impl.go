package service_impl

import (
	"context"
	"github.com/getumen/gogo_crawler/domains/models"
	"github.com/getumen/gogo_crawler/domains/repository"
	"github.com/getumen/gogo_crawler/domains/service"
)

type itemService struct {
	responseRepository repository.ResponseRepository
}

func NewItemService(responseRepository repository.ResponseRepository) service.ItemService {
	return &itemService{responseRepository: responseRepository}
}

func (i *itemService) SaveResponse(ctx context.Context, in <-chan *models.Response) {
	for response := range in {
		_ = i.responseRepository.Save(ctx, response)
	}
}
