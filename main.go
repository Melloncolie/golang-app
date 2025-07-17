package main

import (
	"fmt"
	"golang-app/task"
	"html/template"
	"net/http"
)

func home_page(page http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("tamplates/home_page.html")
	if err != nil {
		fmt.Fprintf(page, err.Error())
	}

	tmpl.Execute(page, task.GetList())
}

func handleRequest() {
	http.HandleFunc("/", home_page)
	http.HandleFunc("/addtask", task.AddToList)
	http.HandleFunc("/deletetask", task.DeleteElement)
	http.HandleFunc("/taskdone", task.TaskDone)
	http.ListenAndServe(":8080", nil)
}
func main() {
	handleRequest()
}
