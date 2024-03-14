package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
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
// j'en étais là, à essayer d'afficher les messages d'erreur côté client.
func (mt *MyTodoList) add(todo Todo) (string, error) {
	existingTodo := false
	if todo.Text == "" {
		return "Pas de texte renseigné !", fmt.Errorf("Pas de texte renseigné !")
	} else {
		for _, t := range myList.todoList {
			if t.Text == todo.Text {
				existingTodo = true
				return "Todo existant !", fmt.Errorf("Todo existant !")
			}
		}
		if existingTodo == false {
			mt.todoList = append(mt.todoList, todo)
		}
	}
	return "", nil
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
	}

	todo := Todo{
		Done: done,
		Text: text,
	}
	// log.Printf("Debug : %v", todo)

	if !todo.Done {
		myList.add(todo)
	} else {
		myList.delete(todo)
	}

	fmt.Println(myList.todoList)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todo)
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(myList.todoList)
}

func main() {
	// handle route using handler function
	fs := http.FileServer(http.Dir("./ihm"))
	http.Handle("/", fs)

	http.HandleFunc("/service", renvoiValeur)
	http.HandleFunc("/todos", getTodos)

	// listen to port
	http.ListenAndServe(":5050", nil)
}
