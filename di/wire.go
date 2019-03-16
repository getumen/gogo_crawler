//+build wireinject

package di

import (
	"github.com/getumen/gogo_crawler/application/usecase"
	"github.com/getumen/gogo_crawler/config"
	"github.com/getumen/gogo_crawler/domains/service_impl"
	"github.com/getumen/gogo_crawler/infras/http"
	"github.com/getumen/gogo_crawler/infras/persistence/mysql"
	"github.com/gomodule/redigo/redis"
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
)

func InitializeCrawler(config *config.Config, db *gorm.DB, redisConn *redis.Pool) (usecase.Crawler, error) {
	wire.Build(
		http.NewHttpClientRepository,
		mysql.NewResponseMysqlRepository,
		mysql.NewRequestMysqlRepository,
		service_impl.NewPoissonProcessRule,
		service_impl.NewCrawlerService,
		usecase.NewCrawler,
	)
	return nil, nil
}
