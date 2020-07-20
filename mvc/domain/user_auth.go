package domain

import "time"

// UserAuth ...
type UserAuth struct {
	UserAuthSeq uint64    `gorm:"primary_key" json:"user_auth_seq"`
	UserSeq     uint64    `gorm:"foreignkey" json:"user_seq"`
	AuthCode    string    `json:"auth_code"`
	AuthName    string    `json:"auth_name"`
	CreatedAt   time.Time `json:"created_at"`
}

// NewUserAuth ...
func NewUserAuth() *UserAuth {
	return &UserAuth{}
}
