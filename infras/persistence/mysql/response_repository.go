package mysql

import (
	"context"
	"github.com/getumen/gogo_crawler/domains/models"
	"github.com/getumen/gogo_crawler/domains/repository"
	"github.com/getumen/gogo_crawler/hashutils"
	"github.com/jinzhu/gorm"
	"time"
)

type responseMysqlRepository struct {
	db *gorm.DB
}

func NewResponseMysqlRepository(db *gorm.DB) repository.ResponseRepository {
	return &responseMysqlRepository{db: db}
}

func newResponseDB(response *models.Response) *responseDB {
	return &responseDB{
		hashutils.HashHTMLBody(response.Body),
		response.UrlString(),
		string(response.Body),
		makeDate(response.CreateAt),
	}

}

func (r *responseMysqlRepository) Save(ctx context.Context, response *models.Response) error {
	resp := newResponseDB(response)
	err := r.db.Save(resp).Error
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (r *responseMysqlRepository) IsExist(ctx context.Context, response *models.Response) (bool, error) {
	resp := newResponseDB(response)
	found := r.db.First(&resp, "response_hash=?", resp.ResponseHash).RecordNotFound()
	return !found, nil
}

func makeDate(d time.Time) int {
	return 10000*d.Year() + 100*int(d.Month()) + d.Day()
}
