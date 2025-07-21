package task

import (
	"encoding/json"
	"fmt"
	"golang-app/config"
	"net/http"
	"os"
	"time"
)

type DataTask struct {
	Tasks    ToDolist
	EditTask *Task
}
type Task struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Isdone      bool      `json:"isdone"`
	Time_start  time.Time `json:"time_start"`
	Time_end    time.Time `json:"time_end"`
}

type ToDolist []Task

func (list *ToDolist) isInList(task Task) bool {
	for _, v := range *list {
		if v.Name == task.Name && v.Description == task.Description {
			fmt.Println("Задача уже существует")
			return true
		}
	}
	return false
}
func (list *ToDolist) GetList() {
	var con config.Config
	con.LoadConfig()
	jsonData, err := os.ReadFile(con.Todolist)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(jsonData, &list)
	if err != nil {
		panic(err)
	}
}

func (list *ToDolist) addToFile() {
	var con config.Config
	con.LoadConfig()
	jsonData, err := json.Marshal(list)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(con.Todolist, jsonData, 0644)
	if err != nil {
		panic(err)
	}
}

func (task *Task) newTask(Name string, Description string) {
	task.Name = Name
	task.Description = Description
	task.Isdone = false
	task.Time_start = time.Now()
	task.Time_end = time.Time{}
}

func (list *ToDolist) AddToList(w http.ResponseWriter, r *http.Request) {
	defer http.Redirect(w, r, "/", http.StatusSeeOther)
	list.GetList()
	name := r.FormValue("taskName")
	description := r.FormValue("description")

	if name == "" || description == "" {
		fmt.Println("Заполнены не все поля")
		http.Error(w, "Заполните все поля", http.StatusBadRequest)
		return
	}

	var task Task
	task.newTask(name, description)
	if list.isInList(task) {
		return
	}
	*list = append(*list, task)
	list.addToFile()
	fmt.Println("Задача ", task, " добавленна в список")
}

func (list *ToDolist) DeleteElement(w http.ResponseWriter, r *http.Request) {
	defer http.Redirect(w, r, "/", http.StatusSeeOther)
	list.GetList()
	name := r.FormValue("taskName")
	for i, v := range *list {
		if v.Name == name {
			*list = append((*list)[:i], (*list)[i+1:]...)
			list.addToFile()
			fmt.Println("Задача ", v, " удалена")
			return
		}
	}

}

func (list *ToDolist) TaskDone(w http.ResponseWriter, r *http.Request) {
	defer http.Redirect(w, r, "/", http.StatusSeeOther)
	list.GetList()
	name := r.FormValue("taskName")
	for i, v := range *list {
		if v.Name == name {
			(*list)[i].Isdone = true
			(*list)[i].Time_end = time.Now()
			list.addToFile()
			fmt.Println("Задача ", v, " отмечена выполненой")
			return
		}
	}
}
func (list *ToDolist) TaskEdit(w http.ResponseWriter, r *http.Request) {
	defer http.Redirect(w, r, "/", http.StatusSeeOther)
	list.GetList()
	oldName := r.FormValue("oldName")
	newName := r.FormValue("newName")
	newDescription := r.FormValue("description")
	var task Task
	task.newTask(newName, newDescription)
	if list.isInList(task) {
		return
	}
	for i, v := range *list {
		if v.Name == oldName {
			if newName != "" {
				(*list)[i].Name = newName
			}
			if newDescription != "" {
				(*list)[i].Description = newDescription
			}
			list.addToFile()
			fmt.Println("Задача", v, "отредактирована")
			break
		}

	}
}
