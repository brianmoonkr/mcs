package domain

import "time"

// CommonCode ...
type CommonCode struct {
	GroupSeq    uint64    `gorm:"primary_key" json:"group_seq"`
	Code        string    `gorm:"primary_key" json:"code"`
	CodeName    string    `json:"code_name"`
	Description string    `json:"description"`
	OrderNum    int       `json:"order_num"`
	UseYN       string    `json:"use_yn"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
