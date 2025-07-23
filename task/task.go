package task

import (
	"encoding/json"
	"fmt"
	"golang-app/config"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Task struct {
	Id          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Isdone      bool      `json:"isdone"`
	Time_start  time.Time `json:"time_start"`
	Time_end    time.Time `json:"time_end"`
}

type ToDolist []Task

func (Tlist *ToDolist) isInList(task Task) bool {
	for _, v := range *Tlist {
		if v.Name == task.Name && v.Description == task.Description {
			fmt.Println("Задача уже существует")
			return true
		}
	}
	return false
}
func (Tlist *ToDolist) GetList() {
	var con config.Config
	con.LoadConfig()
	jsonData, err := os.ReadFile(con.Todolist)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(jsonData, &Tlist)
	if err != nil {
		panic(err)
	}
}

func (Tlist *ToDolist) AddToFile() {
	var con config.Config
	con.LoadConfig()
	jsonData, err := json.Marshal(Tlist)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(con.Todolist, jsonData, 0644)
	if err != nil {
		panic(err)
	}
}

func (Tlist *ToDolist) getId() uint {
	for ID := 0; ID < len(*Tlist); ID++ {
		for i, v := range *Tlist {
			if v.Id == uint(ID) {
				break
			} else if i == len(*Tlist)-1 {
				return uint(ID)
			}
		}
	}
	return uint(len(*Tlist))
}
func (task *Task) newTask(Name string, Description string) {
	var list ToDolist
	list.GetList()
	ID := list.getId()
	task.Id = ID
	task.Name = Name
	task.Description = Description
	task.Isdone = false
	task.Time_start = time.Now()
	task.Time_end = time.Time{}
}

func (Tlist *ToDolist) AddToList(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("taskName")
	description := r.FormValue("description")
	if name == "" || description == "" {
		fmt.Println("Заполнены не все поля")
		http.Error(w, "Заполните все поля", http.StatusBadRequest)
		return
	}
	var task Task
	task.newTask(name, description)
	if Tlist.isInList(task) {
		return
	}
	*Tlist = append(*Tlist, task)
	Tlist.AddToFile()
	fmt.Println("Задача ", task, " добавленна в список")
}

func (Tlist *ToDolist) DeleteElement(w http.ResponseWriter, r *http.Request) {

	taskId, _ := strconv.ParseInt(r.FormValue("taskId"), 10, 10)
	for i, v := range *Tlist {
		if v.Id == uint(taskId) {
			*Tlist = append((*Tlist)[:i], (*Tlist)[i+1:]...)
			Tlist.AddToFile()
			fmt.Println("Задача ", v, " удалена")
			return
		}
	}

}

func (Tlist *ToDolist) TaskDone(w http.ResponseWriter, r *http.Request) {
	Tlist.GetList()
	taskId, _ := strconv.ParseInt(r.FormValue("taskId"), 10, 10)
	for i, v := range *Tlist {
		if v.Id == uint(taskId) {
			(*Tlist)[i].Isdone = true
			(*Tlist)[i].Time_end = time.Now()
			Tlist.AddToFile()
			fmt.Println("Задача ", v, " отмечена выполненой")
			return
		}
	}
}
func (Tlist *ToDolist) TaskEdit(w http.ResponseWriter, r *http.Request) {
	Tlist.GetList()
	oldId, _ := strconv.ParseInt(r.FormValue("oldId"), 10, 10)
	newName := r.FormValue("newName")
	newDescription := r.FormValue("description")
	var task Task
	task.newTask(newName, newDescription)
	if Tlist.isInList(task) {
		return
	}
	for i, v := range *Tlist {
		if v.Id == uint(oldId) {
			if newName != "" {
				(*Tlist)[i].Name = newName
			}
			if newDescription != "" {
				(*Tlist)[i].Description = newDescription
			}
			Tlist.AddToFile()
			fmt.Println("Задача", v, "отредактирована")
			break
		}
	}
}
