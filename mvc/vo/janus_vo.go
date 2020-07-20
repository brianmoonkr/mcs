package vo 

import (

	"github.com/jinzhu/gorm"
	"github.com/teamgrit-lab/cojam/mvc/domain"
)

// JanusVO ...
type JanusVO struct {
	RDBConn *gorm.DB `json:"-"`
	Auth	*domain.Auth	`json:"auth"`	
	AuthCnt         int     `json:"auth_cnt"`
	B2bLive	*domain.B2bLive	`json:"b2b_live"`
}


