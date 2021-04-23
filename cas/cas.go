package cas

import (
	"github.com/go-ini/ini"
	"github.com/imroc/req"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

const (
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.72 Safari/537.36"
	casUrl    = "https://cas.hdu.edu.cn/cas/login"
	contentType ="application/x-www-form-urlencoded"
)

// cas认证会话
type Session struct {
	user    string
	passwd  string
	Request *req.Req
}

func (s *Session) setSession() error {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		return err
	}
	s.user = cfg.Section("cas").Key("user").String()
	s.passwd = cfg.Section("cas").Key("passwd").String()
	s.Request = req.New()

	return nil
}

func (s *Session) casLogin() error {
	lt, execution, err := s.getLoginTicket()
	if err != nil {
		return err
	}
	rsa, err := desEncrypt(s.user + s.passwd + lt)
	if err != nil {
		return err
	}
	param := req.Param{
		"ul":        len(s.user),
		"pl":        len(s.passwd),
		"_eventId":  "submit",
		"execution": execution,
		"lt":        lt,
		"rsa":       rsa,
	}

	resp, err := s.Request.Post(casUrl, req.Header{
		"User-Agent":   userAgent,
		"Content-Type": contentType,
	}, param)
	if err != nil {
		return err
	}

	// login failed
	if strings.Contains(resp.String(), "抱歉！您的请求出现了异常，请稍后再试") {
		return errors.New("cas Bad Request")
	}
	if strings.Contains(resp.String(), "用户名密码错误") {
		return errors.New("cas Account Error")
	}
	return nil
}

func (s *Session) getLoginTicket() (string, string, error) {
	//get 获取lt,execution
	resp, err := s.Request.Get(casUrl, req.Header{
		"User-Agent": userAgent,
	})
	if err != nil {
		return "", "", err
	}
	body := resp.String()
	ltExp := regexp.MustCompile("<input type=\"hidden\" id=\"lt\" name=\"lt\" value=\"(.+)\" />")
	lt := ltExp.FindStringSubmatch(body)[1]
	exExp := regexp.MustCompile("<input type=\"hidden\" name=\"execution\" value=\"(.+)\" />")
	ex := exExp.FindStringSubmatch(body)[1]

	return lt, ex, nil
}
