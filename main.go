package main

import (
	"learn/cas"
	"learn/fanya"
	"log"
	"os"
)

func main() {
	//login cas
	casSession, err := cas.Login()
	if err != nil {
		log.Printf("err:login cas failed,  %+v\n", err)
		return
	}
	//login fanya
	fy := fanya.New(casSession)
	if err := casSession.ServiceLogin(fy); err != nil {
		log.Printf("err:login fanya failed,  %+v\n", err)
		return
	}
	// get courses
	courses, err := fy.GetCourses()
	if err != nil {
		log.Printf("err:get courses failed,  %+v\n", err)
		return
	}
	// get homeworks
	if err := fy.GetHomeworks(courses); err != nil {
		log.Printf("err:get homeworks failed,  %+v\n", err)
		return

	}

	// get html page
	if err := os.Remove("./page/email.html"); err != nil {
		log.Printf("err:remove email.html failed,  %+v\n", err)
		return

	}
	f, err := os.OpenFile("./page/email.html", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := fanya.GetPage(courses, f); err != nil {
		log.Printf("err:get html page failed,  %+v\n", err)
		return
	}
}
