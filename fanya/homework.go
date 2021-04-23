package fanya

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/imroc/req"
	"os"
	"strings"
)
const (
	chaoxingUrl="https://mooc1-2.chaoxing.com"
)
type homework struct {
	name string
	startTime string
	endTime string
	status string
}
func (hw *homework)writeFile(f *os.File)error{
	if hw.status!="待做"{
		return nil
	}

	data:=fmt.Sprintf("%v	开始时间:%v	结束时间:%v	作业状态:%v\n",hw.name,hw.startTime,hw.endTime,hw.status)

	if _, err := f.WriteString(data); err != nil {
		panic(err)
	}
	return nil
}
func(fy *fanya)GetHomeworks(courses []Cousrse)error{
	if err:=fy.getHwUrl(courses);err!=nil{
		return err
	}
	for i,_:=range courses{
		resp,err:=fy.Session.Request.Get(courses[i].hwUrl,req.Header{
			"User-Agent": userAgent,
		})
		if err!=nil{
			return err
		}

		dom,err:=goquery.NewDocumentFromReader(strings.NewReader(resp.String()))
		if err!=nil{
			return err
		}
		if err:=getHwAtr(dom,&courses[i]);err!=nil{
			return err
		}
	}
	return nil
}
// 从作业页面解析html得到作业属性
func getHwAtr(dom *goquery.Document, course *Cousrse)error{


	dom.Find(".titTxt").Each(func(i int, selection *goquery.Selection) {
		hw:=&homework{}
		if value,exist:=selection.Find("a").Attr("title");exist{
			hw.name=value
		}
		// get startTime
		startSel:=selection.Find("span").First()
		startTxt:=startSel.Not(".fl").Text()
		hw.startTime=startTxt[strings.Index(startTxt,"：")+len("："):]

		// get endTime
		endTxt:=startSel.Next().Not(".fl").Text()
		hw.endTime=endTxt[strings.Index(endTxt,"：")+len("："):]

		// get status
		statusTxt:=selection.Find("strong").Text()
		statusTxt=strings.ReplaceAll(statusTxt,"\t","")
		statusTxt=strings.ReplaceAll(statusTxt,"\n","")
		statusTxt=strings.ReplaceAll(statusTxt," ","")
		hw.status=statusTxt
		if hw.status=="待做" {
			course.todo=true
		}
		course.Hw=append(course.Hw,*hw)
	})



	return nil
}
func(fy *fanya)getHwUrl(courses []Cousrse)error{
	for i,course:=range courses {

		resp,err:=fy.Session.Request.Get(course.courseUrl,req.Header{
			"User-Agent": userAgent,
		})
		if err!=nil {
			return err
		}
		dom,err:=goquery.NewDocumentFromReader(strings.NewReader(resp.String()))
		if err!=nil{
			return err
		}
		if val,exist:=dom.Find("a[title=\"作业\"]").Attr("data");exist{
			courses[i].hwUrl=chaoxingUrl+val
		}
	}
	return nil
}
