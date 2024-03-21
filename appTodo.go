package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"webstack/config"
	dt "webstack/data"

	_ "github.com/go-sql-driver/mysql"
)

type MyTodoList struct {
	todoList []dt.Todo
}

var (
	todos  = []dt.Todo{}
	myList = MyTodoList{todos}
)

// créer, modifier, supprimer todo
func (mt *MyTodoList) add(todo dt.Todo) error {
	existingTodo := false
	if todo.Text == "" {
		return fmt.Errorf("pas de texte renseigné")
	} else {
		for _, t := range mt.todoList {
			if t.Text == todo.Text {
				existingTodo = true
				return fmt.Errorf("todo existant")
			}
		}
		if !existingTodo {
			mt.todoList = append(mt.todoList, todo)
			dt.AddTodo(todo)
		}
	}
	return nil
}

func (mt *MyTodoList) delete(todo dt.Todo) error {
	for i, t := range mt.todoList {
		if t.Text == todo.Text {
			// Supprime l'élément visé sans changer l'ordre
			mt.todoList = append(mt.todoList[:i], mt.todoList[i+1:]...)
			dt.DeleteTodo(todo)
			return nil
		}
	}
	return fmt.Errorf("todo '%s' introuvable", todo.Text)
}

func (mt *MyTodoList) modify(todo dt.Todo) error {
	for _, t := range mt.todoList {
		if t.Id == todo.Id {
			t.Text = todo.Text
			dt.ModifyTodo(todo)
			return nil
		}
	}
	return fmt.Errorf("Todo '%s' introuvable", todo.Text)
}

func handleClientRequest(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	action := r.FormValue("action")
	text := r.FormValue("text")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Erreur de conversion", http.StatusBadRequest)
		return
	}

	todo := dt.Todo{
		Id:     id,
		Text:   text,
		Action: action,
	}

	if todo.Action == "add" {
		err := myList.add(todo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else if todo.Action == "delete" {
		err := myList.delete(todo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else if todo.Action == "modify" {
		fmt.Println("Todo à modifier :", todo)
		err := myList.modify(todo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todo)
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	list, err := dt.GetTodos()
	if err != nil {
		log.Fatal("getTodos : ", err)
		return
	}
	myList.todoList = list
	json.NewEncoder(w).Encode(myList.todoList)
}

func main() {
	config.ServConfig = config.ServConfig.GetConfig()
	dir := os.Getenv("DIR")
	if dir == "" {
		dir = "./"
	}

	db, err := dt.OpenDb()
	if err != nil {
		log.Fatal(err)
		return
	}
	config.ServConfig.Db = db
	fmt.Println(config.ServConfig)

	fs := http.FileServer(http.Dir(dir))
	http.Handle("/", fs)

	http.HandleFunc("/service", handleClientRequest)
	http.HandleFunc("/todos", getTodos)

	http.ListenAndServe(config.ServConfig.Port, nil)
}
