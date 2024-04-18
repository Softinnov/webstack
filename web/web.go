package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"webstack/metier/todos"
	"webstack/metier/users"
)

const ERR_AJOUT = "erreur d'ajout de votre tâche"
const ERR_SUPR = "erreur de suppression de votre tâche"
const ERR_MODIF = "erreur de modification de votre tâche"
const ERR_GETDATA = "erreur lors de la récupération des données"
const ERR_ENCOD = "erreur d'encodage json"
const ERR_CONV = "erreur de conversion"

type TodoView struct {
	Id       int    `json:"id"`
	Task     string `json:"task"`
	Priority int    `json:"priority"`
}

func NewTodoView(td todos.Todo) (tdv TodoView) {
	tdv.Id = td.Id
	tdv.Task = todos.GetTask(td.Task)
	tdv.Priority = td.Priority
	return tdv
}

func Todos2TodosView(list []todos.Todo) (displayedList []TodoView) {
	for _, todo := range list {
		displayedList = append(displayedList, NewTodoView(todo))
	}
	return displayedList
}

func encodejson(w http.ResponseWriter, a any) (any, error) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(a)
	if err != nil {
		return nil, fmt.Errorf("%v : %v", ERR_ENCOD, err)
	}
	return a, nil
}

func HandleAddTodo(w http.ResponseWriter, r *http.Request) {
	var user users.User
	text := r.FormValue("task")
	priorityStr := r.FormValue("priority")

	priority, err := strconv.Atoi(priorityStr)
	if err != nil {
		err = fmt.Errorf("%v : %v", ERR_CONV, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tokenStr, err := getTokenString(r, COOKIE_NAME)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user = users.SetEmail(getUserEmail(tokenStr))
	task, err := todos.NewTask(text)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todo, err := todos.Add(task, priority, user)
	if err != nil {
		err = fmt.Errorf("%v : %v", ERR_AJOUT, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	todoview := NewTodoView(todo)
	_, err = encodejson(w, todoview)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func HandleDeleteTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id == 0 {
		err = fmt.Errorf("%v : %v", ERR_CONV, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	todo, err := todos.Delete(id)
	if err != nil {
		err = fmt.Errorf("%v : %v", ERR_SUPR, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	todoview := NewTodoView(todo)
	_, err = encodejson(w, todoview)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func HandleModifyTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	text := r.FormValue("task")
	priorityStr := r.FormValue("priority")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		err = fmt.Errorf("%v : %v", ERR_CONV, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	priority, err := strconv.Atoi(priorityStr)
	if err != nil {
		err = fmt.Errorf("%v : %v", ERR_CONV, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	task, err := todos.NewTask(text)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	todo, err := todos.Modify(task, id, priority)
	if err != nil {
		err = fmt.Errorf("%v : %v", ERR_MODIF, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	todoview := NewTodoView(todo)
	_, err = encodejson(w, todoview)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func HandleGetTodos(w http.ResponseWriter, r *http.Request) {
	var user users.User
	tokenStr, err := getTokenString(r, COOKIE_NAME)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user = users.SetEmail(getUserEmail(tokenStr))
	list, err := todos.Get(user)
	displayedList := Todos2TodosView(list)
	if err != nil {
		http.Error(w, ERR_GETDATA, http.StatusInternalServerError)
		return
	}
	_, err = encodejson(w, displayedList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
