package main

import (
	"fmt"
	"golang-app/config"
	"golang-app/task"
	"html/template"
	"net/http"
)

func home_page(page http.ResponseWriter, r *http.Request) {
	var con config.Config
	con.LoadConfig()

	tmpl, err := template.ParseFiles(con.TemplateFile)
	if err != nil {
		fmt.Fprintf(page, err.Error())
	}
	var list task.ToDolist
	list.GetList()
	editName := r.URL.Query().Get("edit")
	var editTask *task.Task
	for _, t := range list {
		if t.Name == editName {
			editTask = &t
			break
		}
	}
	data := task.DataTask{
		Tasks:    list,
		EditTask: editTask,
	}

	tmpl.Execute(page, data)
}

func handleRequest() {
	var con config.Config
	con.LoadConfig()
	var list task.ToDolist
	fmt.Println("Cервер запущен. Конфигурация ", con)
	http.HandleFunc("/", home_page)
	http.HandleFunc("/addtask", list.AddToList)
	http.HandleFunc("/deletetask", list.DeleteElement)
	http.HandleFunc("/taskdone", list.TaskDone)
	http.HandleFunc("/taskedit", list.TaskEdit)
	http.ListenAndServe(con.Port, nil)
}

func main() {
	handleRequest()
}
