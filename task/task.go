package task

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Task struct {
	Name        string
	Description string
	Isdone      bool
	Time_start  time.Time
	Time_end    time.Time
}

type ToDolist []Task

func GetList() ToDolist {
	jsonData, err := os.ReadFile("ToDoList.json")
	if err != nil {
		panic(err)
	}
	var list ToDolist
	err = json.Unmarshal(jsonData, &list)
	if err != nil {
		panic(err)
	}
	return list
}

func addToFile(list ToDolist) {
	jsonData, err := json.Marshal(list)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("ToDoList.json", jsonData, 0644)
	if err != nil {
		panic(err)
	}
}

func newTask(Name string, Description string) Task {
	task := Task{
		Name:        Name,
		Description: Description,
		Isdone:      false,
		Time_start:  time.Now(),
		Time_end:    time.Time{},
	}
	return task
}

func ViewAllList(w http.ResponseWriter, r *http.Request) {
	list := GetList()
	AllList := ""
	for i, v := range list {
		AllList += fmt.Sprintf("Здача %d \nИмя: %s \nОписание: %s \nНачато %d-%02d-%02d %02d:%02d \n", i+1,
			v.Name, v.Description,
			v.Time_start.Year(), v.Time_start.Month(), v.Time_start.Day(), v.Time_start.Hour(), v.Time_start.Minute())
		if v.Isdone {
			AllList += fmt.Sprintf("Выполненно %d-%02d-%02d %02d:%02d\n",
				v.Time_end.Year(), v.Time_end.Month(), v.Time_end.Day(), v.Time_end.Hour(), v.Time_end.Minute())
		} else {
			AllList += "Не выполненно\n"
		}
	}
	fmt.Fprintf(w, AllList)
}

func AddToList(w http.ResponseWriter, r *http.Request) {
	defer http.Redirect(w, r, "/", http.StatusSeeOther)
	list := GetList()
	name := r.FormValue("taskName")
	description := r.FormValue("description")
	for _, v := range list {
		if v.Name == name && v.Description == description {
			return
		}
	}
	if name == "" || description == "" {
		http.Error(w, "Заполните все поля", http.StatusBadRequest)
		return
	}

	task := newTask(name, description)
	list = append(list, task)
	addToFile(list)
}

func DeleteElement(w http.ResponseWriter, r *http.Request) {
	defer http.Redirect(w, r, "/", http.StatusSeeOther)
	list := GetList()
	name := r.FormValue("taskName")
	for i, v := range list {
		if v.Name == name {
			list = append(list[:i], list[i+1:]...)
			addToFile(list)
			return
		}
	}

}

func TaskDone(w http.ResponseWriter, r *http.Request) {
	defer http.Redirect(w, r, "/", http.StatusSeeOther)
	list := GetList()
	name := r.FormValue("taskName")
	for i, v := range list {
		if v.Name == name {
			list[i].Isdone = true
			list[i].Time_end = time.Now()
			addToFile(list)
			return
		}
	}
}
