package mysql

type responseDB struct {
	ResponseHash string `gorm:"size:256"`
	Url          string `gorm:"size:1024"`
	Body         string
	Date         int
}

func (responseDB) TableName() string {
	return "response"
}