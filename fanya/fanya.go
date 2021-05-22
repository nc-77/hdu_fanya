package fanya

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/imroc/req"
	"github.com/pkg/errors"
	"learn/cas"
	"regexp"
	"strings"
)

const (
	fanyaUrl     = "http://hdu.fanya.chaoxing.com/sso/hdu"
	userAgent    = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.72 Safari/537.36"
	contentType  = "application/x-www-form-urlencoded"
	referValue   = "http://hdu.fanya.chaoxing.com/login/auth"
	authurlValue = "http://hdu.fanya.chaoxing.com/tologin?status=2"
)

type fanya struct {
	Session *cas.Session
}

func New(session *cas.Session) *fanya {
	fy := &fanya{Session: session}
	return fy
}
func (fy *fanya) GetServiceUrl() string {
	return fanyaUrl
}

// 模拟提交表单
func (fy *fanya) Login(body string) error {

	param, err := fy.getParam(body)
	if err != nil {
		return err
	}
	// 取得表单提交url <form action="http://passport2.chaoxing.com/loginfanya?_t=1618921809233" method="post"  id="userLogin">
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return errors.Wrap(err, "获取表单提交url失败")
	}
	postUrl, _ := dom.Find("#userLogin").Attr("action")

	resp, err := fy.Session.Request.Post(string(postUrl), req.Header{
		"User-Agent":   userAgent,
		"Content-Type": contentType,
	}, param)

	if err != nil {
		return err
	}

	if strings.Contains(resp.String(), "验证失败") {
		return errors.New("cas验证通过,泛雅登录失败")
	}
	return nil

}

func (fy *fanya) getParam(body string) (req.Param, error) {
	param := req.Param{
		"fid":     fy.getValue("fid", body),
		"uname":   fy.getValue("uname", body),
		"enc":     fy.getValue("enc", body),
		"refer":   referValue,
		"authurl": authurlValue,
		"time":    fy.getValue("time", body),
	}
	return param, nil
}

// only <input type="hidden" name="fid" value="1001"/>
func (fy *fanya) getValue(name string, body string) string {
	comp := regexp.MustCompile(fmt.Sprintf("<input type=\"hidden\" name=\"%v\" value=\"(.+)\"/>", name))
	value := comp.FindSubmatch([]byte(body))[1]
	return string(value)
}
