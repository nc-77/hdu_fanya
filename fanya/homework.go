package fanya

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/imroc/req"
	"strings"
)

const (
	chaoxingUrl = "https://mooc1-2.chaoxing.com"
)

type homework struct {
	Name      string
	StartTime string
	EndTime   string
	Status    string
}

func (fy *fanya) GetHomeworks(courses []Course) error {
	if err := fy.getHwUrl(courses); err != nil {
		return err
	}
	for i := range courses {
		resp, err := fy.Session.Request.Get(courses[i].hwUrl, req.Header{
			"User-Agent": userAgent,
		})
		if err != nil {
			return err
		}

		dom, err := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))
		if err != nil {
			return err
		}

		if err := getHwAtr(dom, &courses[i]); err != nil {
			return err
		}
	}
	return nil
}

// 从作业页面解析html得到作业属性
func getHwAtr(dom *goquery.Document, course *Course) error {

	dom.Find(".titTxt").Each(func(i int, selection *goquery.Selection) {
		hw := &homework{}
		if value, exist := selection.Find("a").Attr("title"); exist {
			hw.Name = value
		}
		// get startTime
		startSel := selection.Find("span").First()
		startTxt := startSel.Not(".fl").Text()
		hw.StartTime = startTxt[strings.Index(startTxt, "：")+len("："):]

		// get endTime
		endTxt := startSel.Next().Not(".fl").Text()
		hw.EndTime = endTxt[strings.Index(endTxt, "：")+len("："):]

		// get status
		statusTxt := selection.Find("strong").Text()
		statusTxt = strings.ReplaceAll(statusTxt, "\t", "")
		statusTxt = strings.ReplaceAll(statusTxt, "\n", "")
		statusTxt = strings.ReplaceAll(statusTxt, " ", "")
		hw.Status = statusTxt
		if hw.Status == "待做" {
			course.todo = true
		}

		course.Hw = append(course.Hw, *hw)
	})

	return nil
}
func (fy *fanya) getHwUrl(courses []Course) error {
	for i, course := range courses {

		resp, err := fy.Session.Request.Get(course.courseUrl, req.Header{
			"User-Agent": userAgent,
		})
		if err != nil {
			return err
		}
		dom, err := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))
		if err != nil {
			return err
		}
		if val, exist := dom.Find("a[title=\"作业\"]").Attr("data"); exist {
			courses[i].hwUrl = chaoxingUrl + val
		}
	}
	return nil
}
