package redis

import (
	"context"
	"github.com/PuerkitoBio/purell"
	"github.com/getumen/gogo_crawler/domains/models"
	"github.com/getumen/gogo_crawler/domains/repository"
	"github.com/gomodule/redigo/redis"
	"github.com/vmihailenco/msgpack"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// redis namespace
const (
	URL = "url:"
	PQ  = "pq:"
)

const (
	HMSET         = "HMSET"
	ZRANGEBYSCORE = "ZRANGEBYSCORE"
	HGETALL       = "HGETALL"
	EXISTS        = "EXISTS"
	ZADD          = "ZADD"
	ZREM          = "ZREM"
)

type requestRedisRepository struct {
	pool *redis.Pool
}

func NewRequestRedisRepository(pool *redis.Pool) repository.RequestRepository {
	return &requestRedisRepository{pool: pool}
}

func (r *requestRedisRepository) IsExist(ctx context.Context, url string) (bool, error) {
	conn, err := r.pool.GetContext(ctx)
	if err != nil {
		return false, err
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	return redis.Bool(conn.Do(EXISTS, URL+normalizeURL(url)))
}

func (r *requestRedisRepository) FindAllByDomainAndBeforeTimeOrderByNextRequest(
	ctx context.Context,
	namespace string,
	now time.Time,
	offset, limit int) ([]*models.Request, error) {
	conn, err := r.pool.GetContext(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	urls, err := redis.Strings(conn.Do(ZRANGEBYSCORE, PQ+namespace, 0, now.Unix(), "LIMIT", 0, limit))
	if err != nil {
		return nil, err
	}

	if len(urls) > 0 {
		_, err = conn.Do(ZREM, redis.Args{}.Add(PQ + namespace).AddFlat(urls)...)
		if err != nil {
			return nil, err
		}
	}

	m := newRequestMap()
	var wg sync.WaitGroup
	wg.Add(len(urls))
	for _, urlStr := range urls {
		go func(urlStr string) {
			conn, err := r.pool.GetContext(ctx)
			if err != nil {
				return
			}
			defer func() {
				err := conn.Close()
				if err != nil {
					log.Println(err)
				}
			}()

			defer wg.Done()
			var req requestRedis
			v, err := redis.Values(conn.Do(HGETALL, URL+normalizeURL(urlStr)))
			if err != nil {
				return
			}
			if err = redis.ScanStruct(v, &req); err == nil {
				m.store(urlStr, req)
			}
		}(urlStr)
	}
	wg.Wait()
	var res []*models.Request
	for _, urlStr := range urls {
		if v, ok := m.load(urlStr); ok {
			r, err := newRequestFromRedis(&v)
			r.SetNamespace(namespace)
			if err == nil {
				res = append(res, r)
			}
		}
	}
	return res, nil
}

type requestMap struct {
	m *sync.Map
}

func newRequestMap() *requestMap {
	return &requestMap{m: &sync.Map{}}
}

func (r *requestMap) load(key string) (requestRedis, bool) {
	val, ok := r.m.Load(key)
	if !ok {
		return requestRedis{}, false
	}
	return val.(requestRedis), true
}

func (r *requestMap) store(key string, req requestRedis) {
	r.m.Store(key, req)
}

func (r *requestRedisRepository) FindByUrl(ctx context.Context, url string) (*models.Request, error) {
	conn, err := r.pool.GetContext(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	v, err := redis.Values(conn.Do(HGETALL, URL+normalizeURL(url)))
	if err != nil {
		return nil, err
	}
	dst := &requestRedis{}
	err = redis.ScanStruct(v, dst)
	if err != nil {
		return nil, err
	}
	return newRequestFromRedis(dst)
}

func (r *requestRedisRepository) Save(ctx context.Context, request *models.Request) error {
	conn, err := r.pool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	_, err = conn.Do(ZADD, PQ+request.Namespace(), request.NextRequest().Unix(), normalizeURL(request.UrlString()))
	if err != nil {
		return err
	}
	redisRequest, err := newRequestRedis(request)
	if err != nil {
		return err
	}
	_, err = conn.Do(HMSET, redis.Args{}.Add(URL + normalizeURL(request.UrlString())).AddFlat(redisRequest)...)
	return err
}

func newRequestRedis(m *models.Request) (*requestRedis, error) {
	r := &requestRedis{}
	r.Url = m.UrlString()
	r.Method = m.Method()
	r.Body = m.Body()
	if b, err := msgpack.Marshal(m.Cookies()); err == nil {
		r.Cookie = b
	}
	r.JobStatus = m.JobStatus()
	r.NextRequest = m.NextRequest().Unix()
	r.LastRequest = m.LastRequest().Unix()
	if b, err := msgpack.Marshal(m.StatsMap()); err == nil {
		r.Stats = b
		return r, nil
	} else {
		return nil, err
	}
}

func newRequestFromRedis(r *requestRedis) (*models.Request, error) {
	u, err := url.Parse(r.Url)
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
	var s map[string]float64

	if err := msgpack.Unmarshal(r.Stats, &s); err == nil {
		m.SetStatsMap(s)
	}

	return m, nil
}

func normalizeURL(u string) string {
	return purell.MustNormalizeURLString(u, purell.FlagsAllNonGreedy)
}
