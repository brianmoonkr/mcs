package domain

import "time"

// User ...
type User struct {
	UserSeq   uint64    `gorm:"primary_key" json:"user_seq"`
	UserID    string    `json:"user_id"`
	UserName  string    `json:"user_name"`
	NickName  string    `json:"nick_name"`
	UseYN     string    `json:"use_yn"`
	Status    string    `json:"status"`
	AccessDT  time.Time `json:"access_dt"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Users ...
type Users struct {
	User
	RowNum     uint64 `json:"row_num"`
	AuthCode   string `json:"auth_code"`
	Title      string `json:"title"`
	StatusName string `json:"status_name"`
}
