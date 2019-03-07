package mysql

import "time"

type responseDB struct {
	ResponseHash string `gorm:"size:256"`
	Url          string `gorm:"size:1024"`
	Body         string
	Date         int
	CreatedAt    time.Time
}

func (responseDB) TableName() string {
	return "response"
}
