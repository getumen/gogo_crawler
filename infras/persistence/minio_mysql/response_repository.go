package minio_mysql

import (
	"context"
	"github.com/getumen/gogo_crawler/domains/models"
	"github.com/getumen/gogo_crawler/domains/repository"
	"github.com/jinzhu/gorm"
	"github.com/minio/minio-go"
)

type responseMysqlRepository struct {
	db        *gorm.DB
	kvsClient *minio.Client
}

func NewResponseMysqlRepository(db *gorm.DB, kvsClient *minio.Client) repository.ResponseRepository {
	return &responseMysqlRepository{db: db, kvsClient: kvsClient}
}

func (*responseMysqlRepository) Save(ctx context.Context, response *models.Response) error {
	panic("implement me")
}

func (*responseMysqlRepository) IsExist(ctx context.Context, response *models.Response) (bool, error) {
	panic("implement me")
}
