package fanya

import (
	"github.com/imroc/req"
	"html/template"
	"os"
	"regexp"
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
	hrefComp := regexp.MustCompile("<a class=\"zmy_pic\"  target=\"_blank\"  href='(.+)' >")
	coursesUrl := hrefComp.FindAllSubmatch(resp.Bytes(), -1)
	nameComp := regexp.MustCompile("<dt name=\"courseNameHtml\">\n    \t\t\t\t\t\t\t(.+)")
	coursesName := nameComp.FindAllSubmatch(resp.Bytes(), -1)
	courses := make([]Course, len(coursesUrl))
	for i := 0; i < len(coursesUrl); i++ {
		courses[i].Name = string(coursesName[i][1])
		courses[i].courseUrl = string(coursesUrl[i][1])
		courses[i].todo = false
	}

	return courses, nil
}
