package redis

type requestRedis struct {
	Url       string `redis:"url"`
	Method    string `redis:"method"`
	Body      string `redis:"body"`
	Cookie    []byte `redis:"cookie"`
	JobStatus int    `redis:"job_status"`
	// UnixTime
	NextRequest int64 `redis:"next_request" `
	// UnixTime
	LastRequest int64 `redis:"last_request"`
	// stats
	Stats []byte `redis:"stats"`
}
