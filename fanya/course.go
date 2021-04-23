package fanya

import (
	"github.com/imroc/req"
	"os"
	"regexp"
)
type Cousrse struct {
	Name      string
	courseUrl string
	hwUrl     string
	todo		bool
	Hw        []homework
}
const(
	courseUrl="https://hdu.fanya.chaoxing.com/courselist/study"
)
func(course *Cousrse)WriteFile(f *os.File)error{
	if !course.todo {
		return nil
	}
	if _,err:=f.WriteString(course.Name+"\n");err!=nil{
		panic(err)
	}
	for _,hw:=range course.Hw {
		if err:=hw.writeFile(f);err!=nil{
			return err
		}
	}
	return nil
}
func(fy *fanya)GetCourses ()([]Cousrse,error){

	resp,err:=fy.Session.Request.Get(courseUrl,req.Header{
		"User-Agent": userAgent,
	})
	if err!=nil{
		return nil, err
	}
	hrefComp:=regexp.MustCompile("<a class=\"zmy_pic\"  target=\"_blank\"  href='(.+)' >")
	coursesUrl:=hrefComp.FindAllSubmatch(resp.Bytes(),-1)
	nameComp:=regexp.MustCompile("<dt name=\"courseNameHtml\">\n    \t\t\t\t\t\t\t(.+)")
	coursesName:=nameComp.FindAllSubmatch(resp.Bytes(),-1)
	courses:=make([]Cousrse,len(coursesUrl))
	for i:=0;i<len(coursesUrl);i++{
		courses[i].Name=string(coursesName[i][1])
		courses[i].courseUrl =string(coursesUrl[i][1])
		courses[i].todo=false
	}

	return courses,nil
}


