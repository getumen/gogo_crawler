package mysql

import (
	"context"
	"crypto/sha512"
	"fmt"
	"github.com/getumen/gogo_crawler/domains/models"
	"github.com/getumen/gogo_crawler/domains/repository"
	"github.com/jinzhu/gorm"
	"github.com/vmihailenco/msgpack"
	"io"
	"net/http"
	"net/url"
	"time"
)

type requestMysqlRepository struct {
	db *gorm.DB
}

func NewRequestMysqlRepository(db *gorm.DB) repository.RequestRepository {
	return &requestMysqlRepository{db: db}
}

func (r *requestMysqlRepository) IsExist(ctx context.Context, url string) (bool, error) {
	found := !r.db.First(&requestDB{}, "url_hash = ?", hash(url)).RecordNotFound()
	return found, nil
}

func (r *requestMysqlRepository) FindAllByDomainAndBeforeTimeOrderByNextRequest(ctx context.Context,
	domain string,
	now time.Time,
	offset,
	limit int) ([]*models.Request, error) {
	reqs := []requestDB{}
	r.db.Limit(limit).Offset(offset).Where(
		"next_request < ? AND namespace = ?",
		now.Unix(),
		domain).Order("next_request asc").Find(&reqs)
	requests := []*models.Request{}
	urlHashList := []string{}
	for _, req := range reqs {
		request, err := newRequestFromDB(&req)
		if err != nil {
			continue
		}
		requests = append(requests, request)
		urlHashList = append(urlHashList, req.UrlHash)
	}
	r.db.Where("url_hash in (?)", urlHashList).Delete(requestDB{})
	return requests, nil
}

func (r *requestMysqlRepository) FindByUrl(ctx context.Context, url string) (*models.Request, error) {
	key := hash(url)
	req := &requestDB{}
	r.db.Where(&requestDB{UrlHash: key}).First(req)
	return newRequestFromDB(req)
}

func (r *requestMysqlRepository) Save(ctx context.Context, req *models.Request) error {
	request, err := newRequestDB(req)
	if err != nil {
		return err
	}
	err = r.db.Save(request).Error
	if err != nil {
		return err
	} else {
		return nil
	}
}

func hash(url string) string {
	h := sha512.New()
	_, _ = io.WriteString(h, url)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func newRequestDB(m *models.Request) (*requestDB, error) {
	r := &requestDB{}
	r.Url = m.UrlString()
	r.UrlHash = hash(m.UrlString())
	r.Method = m.Method()
	r.Body = m.Body()
	if b, err := msgpack.Marshal(m.Cookies()); err == nil {
		r.Cookie = b
	}
	r.JobStatus = m.JobStatus()
	r.NextRequest = m.NextRequest().Unix()
	r.LastRequest = m.LastRequest().Unix()
	r.Namespace = m.Namespace()
	if b, err := msgpack.Marshal(m.StatsMap()); err == nil {
		r.Stats = b
		return r, nil
	} else {
		return nil, err
	}
}

func newRequestFromDB(r *requestDB) (*models.Request, error) {
	u, err := url.ParseRequestURI(r.Url)
	if err != nil {
		return nil, err
	}

	m := models.NewRequest(u, r.Method, r.Body)
	var c []*http.Cookie

	if err := msgpack.Unmarshal(r.Cookie, &c); err == nil {
		m.SetCookies(c)
	}

	m.SetJobStatus(r.JobStatus)
	m.SetNextRequest(time.Unix(r.NextRequest, 0))
	m.SetLastRequest(time.Unix(r.LastRequest, 0))
	m.SetNamespace(r.Namespace)
	var s map[string]float64

	if err := msgpack.Unmarshal(r.Stats, &s); err == nil {
		m.SetStatsMap(s)
	}

	return m, nil
}
