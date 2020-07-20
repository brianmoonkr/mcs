package dao

import (
	"github.com/jinzhu/gorm"
	"github.com/teamgrit-lab/cojam/mvc/domain"
	"github.com/teamgrit-lab/cojam/mvc/vo"
)

// NewLiveDAO ...
func NewJanusDAO() JanusDAO {
	return &janusDAO{}
}

// JanusDAO ...
type JanusDAO interface {
	CountServiceAuth(*vo.JanusVO) (int, error)
	SelectAuthByKey(*vo.JanusVO) (*domain.Auth, error)
	CreateLive(*vo.JanusVO) error
	SelectLiveInfo(*vo.JanusVO) (*domain.B2bLive, error)
	SelectLiveInfoByVideoroom(*vo.JanusVO) (*domain.B2bLive, error)
	DeleteLive(*vo.JanusVO) error
	SelectLiveRoomListByAccount(*vo.JanusVO) ([]*domain.MediaChannel, error)
	UpdateLive(*vo.JanusVO) error
}

type janusDAO struct{}

func (janusDAO) CountServiceAuth(vo *vo.JanusVO) (cnt int, err error) {
	tx := vo.RDBConn.Table("b2b_auth a").
		Select(`
			count(*)
		`).
		Where("a.key = ?", vo.Auth.Key).
		Where("a.pwd = ?", vo.Auth.Password).
		Where("a.ip = ?", vo.Auth.RemoteAddr)

	err = tx.Count(&cnt).Error

	return
}

func (janusDAO) SelectAuthByKey(vo *vo.JanusVO) (auth *domain.Auth, err error) {
	auth = &domain.Auth{}

	err = vo.RDBConn.Table("b2b_auth a").
		Where("a.key = ?", vo.Auth.Key).
		Find(auth).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	return auth, nil
}

func (janusDAO) CreateLive(vo *vo.JanusVO) error {
	return vo.RDBConn.Create(vo.B2bLive).Error
}

func (janusDAO) SelectLiveInfo(vo *vo.JanusVO) (live *domain.B2bLive, err error) {
	live = &domain.B2bLive{}

	err = vo.RDBConn.Table("b2b_live a").
		Where("a.live_seq = ?", vo.B2bLive.LiveSeq).
		Find(live).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	return live, nil
}

func (janusDAO) SelectLiveInfoByVideoroom(vo *vo.JanusVO) (live *domain.B2bLive, err error) {
	live = &domain.B2bLive{}

	err = vo.RDBConn.Table("b2b_live a").
		Where("a.videoroom = ?", vo.B2bLive.Videoroom).
		Find(live).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	return live, nil
}

func (janusDAO) DeleteLive(vo *vo.JanusVO) error {
	/*live := &domain.B2bLive{
		Status: "9000",
	}*/
	
	result := vo.RDBConn.
		Where("live_seq = ?", vo.B2bLive.LiveSeq).
		Find(vo.B2bLive).
		//Model(vo.B2bLive).
		Update("Status", "9000")

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (janusDAO) SelectLiveRoomListByAccount(vo *vo.JanusVO) (b2bLives []*domain.MediaChannel, err error) {
	err = vo.RDBConn.Table("b2b_live a").
		Select(`
			 a.live_seq
			,a.room_name
			,a.access_point
			,a.videoroom
			,a.textroom
			,a.description
			,a.record
		`).
		Joins("inner join b2b_auth b on b.seq = a.service_seq").
		Where("b.key = ?", vo.Auth.Key).
		Where("a.status = ?", vo.B2bLive.Status).
		Order("a.created_at desc").
		Scan(&b2bLives).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	return
}

func (janusDAO) UpdateLive(vo *vo.JanusVO) error {
	live := &domain.B2bLive{
		Status: "1000",
	}

	result := vo.RDBConn.
		Model(vo.B2bLive).
		Update(live)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
