package main

import (
	"github.com/pkg/errors"
	"learn/cas"
	"learn/fanya"
	"os"
)

func main() {
	//login cas
	casSession, err := cas.Login()
	if err != nil {
		errors.Wrap(err, "login cas failed")
	}
	//login fanya
	fy := fanya.New(casSession)
	if err := casSession.ServiceLogin(fy); err != nil {
		errors.Wrap(err, "login fanya failed")
	}
	// get courses
	courses,err:=fy.GetCourses()
	if err!=nil{
		errors.Wrap(err,"get courses failed")
	}
	// get homeworks
	if err:=fy.GetHomeworks(courses);err!=nil{
		errors.Wrap(err,"get homeworks failed")
	}
	// output
	if err:=os.Remove("todo.txt");err!=nil{
		errors.Wrap(err,"remove todo.txt failed")
	}
	f, err := os.OpenFile("todo.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	for _,course:=range courses{
		if err:=course.WriteFile(f);err!=nil{
			errors.Wrap(err,"output todo.txt failed")
		}
	}
}
