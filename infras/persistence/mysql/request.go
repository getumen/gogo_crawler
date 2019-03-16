package mysql

type requestDB struct {
	Url       string `gorm:"size:256"`
	UrlHash   string `gorm:"primary_key"`
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

func (requestDB) TableName() string {
	return "request"
}
