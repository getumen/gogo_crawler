package cassandra

type request struct {
	Url       string
	Method    string
	Body      []byte
	Cookie    []byte
	JobStatus int
	// UnixTime
	NextRequest int64
	// UnixTime
	LastRequest int64
	// stats
	Stats     []byte
	Namespace string
}
