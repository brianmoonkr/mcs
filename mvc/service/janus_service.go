package service

import (
	"github.com/teamgrit-lab/cojam/mvc/dao"
	"github.com/teamgrit-lab/cojam/mvc/vo"
	"github.com/teamgrit-lab/cojam/mvc/domain"
)

// NewLiveService ...
func NewJanusService() JanusService {
	return &janusService{}
}

type janusService struct {
	janusDAO            dao.JanusDAO
	
}

type JanusService interface {
	SetService(map[string]interface{})
	GetServiceAuthInfo(*vo.JanusVO) (bool, error)
	RegisterLive(*vo.JanusVO) error
	GetLiveInfo(*vo.JanusVO) (*domain.B2bLive, error)
	GetAuthByKey(*vo.JanusVO) (*domain.Auth, error)
	GetLiveRoomListByAccount(*vo.JanusVO) ([]*domain.MediaChannel, error)
	DeleteLiveInfo(*vo.JanusVO) error
	SetLiveInfo(*vo.JanusVO) error
}

func (s *janusService) SetService(di map[string]interface{}) {
	s.janusDAO = di["janusDAO"].(dao.JanusDAO)
	
}

func (s *janusService) GetServiceAuthInfo(vo *vo.JanusVO) (bool, error) {

	cnt, _ := s.janusDAO.CountServiceAuth(vo)

	if cnt < 1 {
		return false, nil
	}

	return true, nil
}

func (s *janusService) GetLiveInfo(vo *vo.JanusVO) (live *domain.B2bLive, err error) {
	live = &domain.B2bLive{}
	
	if vo.B2bLive.LiveSeq != "" {
		live, err = s.janusDAO.SelectLiveInfo(vo)
	} else if vo.B2bLive.Videoroom != "" {
		live, err = s.janusDAO.SelectLiveInfoByVideoroom(vo)
	}
	if err != nil {
		return nil, err
	}

	return live, nil
}


func (s *janusService) RegisterLive(vo *vo.JanusVO) error {
	
	auth, err := s.janusDAO.SelectAuthByKey(vo)
	if err != nil {
		return err
	}

	vo.B2bLive.ServiceSeq = auth.Seq

	err = s.janusDAO.CreateLive(vo)
	if err != nil {
		return err
	}

	return nil
}

func (s *janusService) GetAuthByKey(vo *vo.JanusVO) (auth *domain.Auth, err error) {
	
	auth, err = s.janusDAO.SelectAuthByKey(vo)
	if err != nil {
		return nil, err
	}

	return auth, nil
}

func (s *janusService) GetLiveRoomListByAccount(vo *vo.JanusVO) (b2bLives []*domain.MediaChannel, err error) {
	b2bLives, err = s.janusDAO.SelectLiveRoomListByAccount(vo)
	
	if err != nil {
		return nil, err
	}

	return b2bLives, nil
}

func (s *janusService) DeleteLiveInfo(vo *vo.JanusVO) error {
	err := s.janusDAO.DeleteLive(vo)
	if err != nil {
		return err
	}

	return nil
}

func (s *janusService) SetLiveInfo(vo *vo.JanusVO) error {
	err := s.janusDAO.UpdateLive(vo)
	if err != nil {
		return err
	}

	return nil
}