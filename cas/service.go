package cas

import (
	"github.com/imroc/req"
)

type Service interface {
	GetServiceUrl() string
	Login(string) error
}

func Login() (*Session, error) {
	s := &Session{}
	if err := s.setSession(); err != nil {
		return s, err
	}
	if err := s.casLogin(); err != nil {
		return s, err
	}
	return s, nil
}

func (s *Session) ServiceLogin(svc Service) error {
	resp, err := s.Request.Get(svc.GetServiceUrl(), req.Header{
		"User-Agent": userAgent,
	})
	if err != nil {
		return err
	}

	if err := svc.Login(resp.String()); err != nil {
		return err
	}
	return nil
}
