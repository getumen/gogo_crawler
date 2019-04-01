package cassandra

import (
	"context"
	"fmt"
	"github.com/getumen/gogo_crawler/domains/models"
	"github.com/getumen/gogo_crawler/domains/repository"
	"github.com/gocql/gocql"
	"github.com/vmihailenco/msgpack"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type requestRepository struct {
	session *gocql.Session
}

func NewRequestCassandraRepository(session *gocql.Session) repository.RequestRepository {
	return &requestRepository{session: session}
}

func (r *requestRepository) IsExist(ctx context.Context, namespace, url string) (bool, error) {
	_, err := r.FindByUrl(ctx, namespace, url)
	if err == gocql.ErrNotFound {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func (r *requestRepository) FindAllByDomainAndBeforeTimeOrderByNextRequest(
	ctx context.Context,
	namespace string,
	now time.Time,
	offset,
	limit int) ([]*models.Request, error) {
	if offset > 0 {
		panic("offset is not supported")
	}

	var reqs []*models.Request

	iter := r.session.Query("SELECT *"+
		" FROM request_pq"+
		" WHERE namespace = ?"+
		" AND next_request < ?"+
		" ORDER BY next_request ASC, url ASC"+
		" LIMIT ?", namespace, now.Unix(), limit).Iter()

	//log.Printf("request_pq = %d", iter.NumRows())

	batch := r.session.NewBatch(gocql.LoggedBatch)
	var urls []string

	sliceMap, err := iter.SliceMap()

	if err != nil {
		return nil, err
	}

	for _, m := range sliceMap {
		urls = append(urls, m["url"].(string))
		batch.Query("DELETE FROM request_pq WHERE namespace = ? AND next_request = ? AND url = ?",
			namespace, m["next_request"], m["url"])
	}

	qu := make([]string, len(urls))
	for i := 0; i < len(qu); i++ {
		qu[i] = "?"
	}
	args := make([]interface{}, len(urls)+1)
	args[0] = namespace
	for i := 0; i < len(urls); i++ {
		args[i+1] = urls[i]
	}

	iter = r.session.Query(fmt.Sprintf("SELECT *"+
		" FROM request"+
		" WHERE namespace = ?"+
		" AND url in (%s)", strings.Join(qu, ",")), args...).Iter()

	sliceMap, err = iter.SliceMap()

	if err != nil {
		return nil, err
	}

	for _, m := range sliceMap {
		req, err := newRequestFromDB(&request{
			Url:         m["url"].(string),
			Method:      m["method"].(string),
			Body:        m["body"].([]byte),
			Cookie:      m["cookie"].([]byte),
			JobStatus:   m["job_status"].(int),
			NextRequest: m["next_request"].(int64),
			LastRequest: m["last_request"].(int64),
			Stats:       m["stats"].([]byte),
			Namespace:   m["namespace"].(string),
		})

		if err != nil {
			log.Println(err)
			continue
		}
		reqs = append(reqs, req)
	}
	err = r.session.ExecuteBatch(batch)
	if err != nil {
		return nil, err
	}

	return reqs, nil
}

func (r *requestRepository) FindByUrl(ctx context.Context, namespace, url string) (*models.Request, error) {
	m := map[string]interface{}{}
	iter := r.session.Query("SELECT * FROM request WHERE namespace = ? AND url = ?", namespace, url).Iter()
	for iter.MapScan(m) {
		return newRequestFromDB(&request{
			Url:         m["url"].(string),
			Method:      m["method"].(string),
			Body:        m["body"].([]byte),
			Cookie:      m["cookie"].([]byte),
			JobStatus:   m["job_status"].(int),
			NextRequest: m["next_request"].(int64),
			LastRequest: m["last_request"].(int64),
			Stats:       m["stats"].([]byte),
			Namespace:   m["namespace"].(string),
		})
	}
	return nil, gocql.ErrNotFound
}

func (r *requestRepository) Save(ctx context.Context, m *models.Request) error {
	req, err := newRequestDB(m)
	if err != nil {
		return err
	}
	b := r.session.NewBatch(gocql.LoggedBatch)
	b.Query("INSERT INTO request_pq(namespace, url, next_request) VALUES(?, ?, ?)",
		req.Namespace, req.Url, req.NextRequest)
	b.Query("INSERT INTO request(namespace, url, method, body, job_status, next_request, last_request, cookie, stats)"+
		" VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		req.Namespace, req.Url, req.Method, req.Body, req.JobStatus, req.NextRequest, req.LastRequest, req.Cookie, req.Stats)
	return r.session.ExecuteBatch(b)
}

func newRequestDB(m *models.Request) (*request, error) {
	r := &request{}
	r.Url = m.UrlString()
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

func newRequestFromDB(r *request) (*models.Request, error) {
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
