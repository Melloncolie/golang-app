package main

import (
	"fmt"
	"golang-app/config"
	"golang-app/task"
	"golang-app/user"
	"html/template"
	"net/http"
	"strconv"
)

func home_page(page http.ResponseWriter, r *http.Request) {
	var con config.Config
	con.LoadConfig()

	tmpl, err := template.ParseFiles(con.HomePage)
	if err != nil {
		fmt.Fprintf(page, err.Error())
	}

	var Uslist user.UserList
	Uslist.GetList()

	var us user.User
	for _, v := range Uslist {
		if v.Success {
			us = v
			break
		}
	}
	var Tlist task.ToDolist
	Tlist = us.List
	var editTask *task.Task
	if editStr := r.URL.Query().Get("edit"); editStr != "" {
		editId, _ := strconv.ParseInt(r.URL.Query().Get("edit"), 10, 10)
		for _, t := range Tlist {
			if t.Id == uint(editId) {
				editTask = &t
				break
			}
		}
	}

	if us.Success {
		data := user.Data{
			User_name: us,
			EditTask:  editTask,
		}
		tmpl.Execute(page, data)
	} else {
		http.Redirect(page, r, "/log", http.StatusSeeOther)
	}
}

func log_page(page http.ResponseWriter, r *http.Request) {
	var con config.Config
	con.LoadConfig()
	tmpl, err := template.ParseFiles(con.LogPage)
	if err != nil {
		fmt.Fprintf(page, err.Error())
	}
	tmpl.Execute(page, nil)
}

func reg_page(page http.ResponseWriter, r *http.Request) {
	var con config.Config
	con.LoadConfig()
	tmpl, err := template.ParseFiles(con.RegPage)
	if err != nil {
		fmt.Fprintf(page, err.Error())
	}
	tmpl.Execute(page, nil)
}

func handleRequest() {
	var con config.Config
	con.LoadConfig()
	var Uslist user.UserList
	fmt.Println("Cервер запущен. Конфигурация ", con)
	http.HandleFunc("/log", log_page)
	http.HandleFunc("/log/oldUser", Uslist.Log_in)
	http.HandleFunc("/reg", reg_page)
	http.HandleFunc("/reg/newUser", Uslist.Reg)
	http.HandleFunc("/", home_page)
	http.HandleFunc("/addtask", Uslist.AddToList)
	http.HandleFunc("/deletetask", Uslist.DeleteElement)
	http.HandleFunc("/taskdone", Uslist.TaskDone)
	http.HandleFunc("/taskedit", Uslist.TaskEdit)
	http.HandleFunc("/exit", Uslist.Exit)
	http.ListenAndServe(con.Port, nil)
}

func main() {
	handleRequest()
}
