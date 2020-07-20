package domain

import "time"

// CommonCodeGroup ...
type CommonCodeGroup struct {
	GroupSeq  uint64    `gorm:"primary_key" json:"group_seq"`
	GroupName string    `json:"group_name"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
