package user

import (
	"encoding/json"
	"fmt"
	"golang-app/config"
	"golang-app/task"
	"net/http"
	"os"
	"strings"
)

type Data struct {
	User_name User
	EditTask  *task.Task
}

type User struct {
	Login    string        `json:"login"`
	Password string        `json:"password"`
	Email    string        `json:"email"`
	Success  bool          `json:"success"`
	List     task.ToDolist `json:"list"`
}

type UserList []User

func (Uslist *UserList) checkAll(login string, password string, checkpassword string, email string) bool {

	if login == "" || password == "" || checkpassword == "" || email == "" {
		fmt.Println("Заполнены не все поля")
		return false
	}
	if password != checkpassword {
		fmt.Println("Пароли не совпадают")
		return false
	}
	if !strings.Contains(email, "@") {
		fmt.Println("Email не содержит @")
		return false
	}
	Uslist.GetList()
	for _, v := range *Uslist {
		if v.Email == email {
			fmt.Println("Email уже зарегестрирован")
			return false
		}
		if v.Login == login {
			fmt.Println("Пользователь с таким именем уже зарегестрирован")
			return false
		}
	}
	return true
}
func (Uslist UserList) addToFileHelper(con config.Config) {
	jsonData, err := json.Marshal(Uslist)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(con.Userlist, jsonData, 0644)
	if err != nil {
		panic(err)
	}
}
func (user *User) addToFile() {
	var con config.Config
	con.LoadConfig()
	var Uslist UserList
	Uslist.GetList()
	for i, v := range Uslist {
		if v.Login == user.Login {
			Uslist[i] = *user
			Uslist.addToFileHelper(con)
			return
		}
	}
	Uslist = append(Uslist, *user)
	Uslist.addToFileHelper(con)
}

func (Uslist *UserList) GetList() {
	var con config.Config
	con.LoadConfig()
	jsonData, err := os.ReadFile(con.Userlist)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(jsonData, &Uslist)
	if err != nil {
		panic(err)
	}
}

func (user *User) newUser(login string, password string, email string) {
	var Uslist UserList
	var Tlist task.ToDolist
	Uslist.GetList()
	user.Login = login
	user.Password = password
	user.Email = email
	user.Success = true
	user.List = Tlist
}

func (Uslist *UserList) Reg(w http.ResponseWriter, r *http.Request) {
	defer http.Redirect(w, r, "/", http.StatusSeeOther)
	Uslist.GetList()
	login := r.FormValue("login")
	password := r.FormValue("password")
	checkpassword := r.FormValue("checkpassword")
	email := r.FormValue("email")
	if Uslist.checkAll(login, password, checkpassword, email) {
		var user User
		user.newUser(login, password, email)
		*Uslist = append(*Uslist, user)
		user.addToFile()
		fmt.Println("Аккаунт ", user, " успешно зарегистрирован")
		return
	}
	http.Redirect(w, r, "/reg", http.StatusSeeOther)
}

func (Uslist *UserList) Exit(w http.ResponseWriter, r *http.Request) {
	tr := r.FormValue("tr")
	if tr == "1" {
		defer http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	for i, v := range *Uslist {
		if v.Success {
			(*Uslist)[i].Success = false
			(*Uslist)[i].addToFile()
			fmt.Println("Выход из аккаунта ", v)
		}
	}

}
func (Uslist *UserList) Log_in(w http.ResponseWriter, r *http.Request) {
	defer http.Redirect(w, r, "/", http.StatusSeeOther)
	Uslist.GetList()
	Uslist.Exit(w, r)
	login := r.FormValue("login")
	password := r.FormValue("password")
	for i, v := range *Uslist {
		if v.Login == login && v.Password == password {
			(*Uslist)[i].Success = true
			Tlist := v.List
			Tlist.AddToFile()
			(*Uslist)[i].addToFile()
			fmt.Println("Вход в аккаунт ", v)
			return
		}
	}
	fmt.Println("Неверный логин или пароль")
	http.Redirect(w, r, "/log", http.StatusSeeOther)
}

func (Uslist *UserList) AddToList(w http.ResponseWriter, r *http.Request) {
	defer http.Redirect(w, r, "/", http.StatusSeeOther)
	var Tlist task.ToDolist
	Tlist.GetList()

	for i, v := range *Uslist {
		if v.Success {
			Tlist.AddToList(w, r)
			Tlist.GetList()
			(*Uslist)[i].List = Tlist
			(*Uslist)[i].addToFile()
			return
		}
	}
}

func (Uslist *UserList) DeleteElement(w http.ResponseWriter, r *http.Request) {
	defer http.Redirect(w, r, "/", http.StatusSeeOther)
	var Tlist task.ToDolist
	Tlist.GetList()

	for i, v := range *Uslist {
		if v.Success {
			Tlist.DeleteElement(w, r)
			Tlist.GetList()
			(*Uslist)[i].List = Tlist
			(*Uslist)[i].addToFile()
			return
		}
	}
}

func (Uslist *UserList) TaskDone(w http.ResponseWriter, r *http.Request) {
	defer http.Redirect(w, r, "/", http.StatusSeeOther)
	var Tlist task.ToDolist
	Tlist.GetList()
	for i, v := range *Uslist {
		if v.Success {
			Tlist.TaskDone(w, r)
			Tlist.GetList()
			(*Uslist)[i].List = Tlist
			(*Uslist)[i].addToFile()
			return
		}
	}
}

func (Uslist *UserList) TaskEdit(w http.ResponseWriter, r *http.Request) {
	defer http.Redirect(w, r, "/", http.StatusSeeOther)
	var Tlist task.ToDolist
	Tlist.GetList()
	for i, v := range *Uslist {
		if v.Success {
			Tlist.TaskEdit(w, r)
			Tlist.GetList()
			(*Uslist)[i].List = Tlist
			(*Uslist)[i].addToFile()
			return
		}
	}
}
