package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type Todo struct {
	Done bool   `json:"done"`
	Text string `json:"text"`
}

type MyTodoList struct {
	todoList []Todo
}

var todos = []Todo{}
var myList = MyTodoList{todos}

// créer, modifier, supprimer todo
func (mt *MyTodoList) add(todo Todo) error {
	existingTodo := false
	if todo.Text == "" {
		return fmt.Errorf("Pas de texte renseigné !")
	} else {
		for _, t := range myList.todoList {
			if t.Text == todo.Text {
				existingTodo = true
				return fmt.Errorf("Todo existant !")
			}
		}
		if !existingTodo {
			mt.todoList = append(mt.todoList, todo)
		}
	}
	return nil
}

func (mt *MyTodoList) delete(todo Todo) error {
	for i, t := range mt.todoList {
		if t.Text == todo.Text {
			//Supprime l'élément visé sans changer l'ordre
			mt.todoList = append(mt.todoList[:i], mt.todoList[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Todo '%s' introuvable", todo.Text)
}

func (mt *MyTodoList) modif(oldText, newText string) error {
	for i, t := range mt.todoList {
		if t.Text == oldText {
			mt.todoList[i].Text = newText
			return nil
		}
	}
	return fmt.Errorf("Todo '%s' introuvable", oldText)
}

func renvoiValeur(w http.ResponseWriter, r *http.Request) {
	doneStr := r.FormValue("check")
	text := r.FormValue("text")

	done, err := strconv.ParseBool(doneStr)

	if err != nil {
		log.Println("Error in conversion", err)
		http.Error(w, "Erreur de conversion", http.StatusBadRequest)
		return
	}

	todo := Todo{
		Done: done,
		Text: text,
	}
	// log.Printf("Debug : %v", todo)

	if !todo.Done {
		err := myList.add(todo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		err := myList.delete(todo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	// fmt.Println(myList.todoList)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todo)
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(myList.todoList)
}

func getDb() {
	db, err := sql.Open("mysql", "adminUser:adminPassword@127.0.0.1:3306/todos")
	if err != nil {
		fmt.Println(err.Error())
	}
	// defer db.Close()
	fmt.Println(db)

	err = db.Ping()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func main() {
	// handle route using handler function
	fs := http.FileServer(http.Dir("./ihm"))
	http.Handle("/", fs)

	http.HandleFunc("/service", renvoiValeur)
	http.HandleFunc("/todos", getTodos)

	// listen to port
	http.ListenAndServe(":5050", nil)

	getDb()
}
