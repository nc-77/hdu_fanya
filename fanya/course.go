package fanya

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/imroc/req"
	"github.com/pkg/errors"
	"html/template"
	"os"
	"strings"
)

type Course struct {
	Name      string
	courseUrl string
	hwUrl     string
	todo      bool
	Hw        []homework
}
type Data struct {
	Courses []Course
}

const (
	courseUrl = "https://hdu.fanya.chaoxing.com/courselist/study"
	cxurl     = "https://hdu.fanya.chaoxing.com"
)

func GetPage(courses []Course, file *os.File) error {
	tmpl, err := template.ParseFiles("./page/template.html")
	if err != nil {
		return err
	}
	data := Data{}
	for _, course := range courses {
		if course.todo {
			todoHw := make([]homework, 0)
			for _, hw := range course.Hw {
				if hw.Status == "待做" {
					todoHw = append(todoHw, hw)
				}
			}
			course.Hw = todoHw
			data.Courses = append(data.Courses, course)
		}
	}

	if err := tmpl.Execute(file, data); err != nil {
		return err
	}
	return nil
}

func (fy *fanya) GetCourses() ([]Course, error) {

	resp, err := fy.Session.Request.Get(courseUrl, req.Header{
		"User-Agent": userAgent,
	})
	if err != nil {
		return nil, err
	}

	dom, err := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))
	if err != nil {
		return nil, errors.WithMessage(err, "courses page to dom failed")
	}
	courses := make([]Course, 0)
	dom.Find("li.zmy_item").Each(func(i int, selection *goquery.Selection) {
		url, _ := selection.Find("a.zmy_pic").Attr("href")
		courseName := selection.Find("dt[name=courseNameHtml]").Text()
		courseName = strings.ReplaceAll(courseName, "\t", "")
		courseName = strings.ReplaceAll(courseName, " ", "")
		courseName = strings.ReplaceAll(courseName, "\n", "")
		courses = append(courses, Course{
			Name:      courseName,
			courseUrl: cxurl + url,
			hwUrl:     "",
			todo:      false,
			Hw:        nil,
		})
	})

	return courses, nil
}
