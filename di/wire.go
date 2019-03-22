//+build wireinject

package di

import (
	"github.com/getumen/gogo_crawler/application/usecase"
	"github.com/getumen/gogo_crawler/config"
	"github.com/getumen/gogo_crawler/domains/service_impl"
	"github.com/getumen/gogo_crawler/infras/http"
	"github.com/getumen/gogo_crawler/infras/persistence/cassandra"
	"github.com/getumen/gogo_crawler/infras/persistence/mysql"
	"github.com/gocql/gocql"
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
)

func InitializeCrawler(config *config.Config, db *gorm.DB, sess *gocql.Session) (usecase.Crawler, error) {
	wire.Build(
		http.NewHttpClientRepository,
		cassandra.NewRequestCassandraRepository,
		mysql.NewResponseMysqlRepository,
		service_impl.NewPoissonProcessRule,
		service_impl.NewCrawlerService,
		usecase.NewCrawler,
	)
	return nil, nil
}
