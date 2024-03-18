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

type Config struct {
	Port string
	Db   *sql.DB
}

type Todo struct {
	Done bool   `json:"done"`
	Text string `json:"text"`
}

type MyTodoList struct {
	todoList []Todo
}

var config = Config{
	Port: ":5050",
}

var todos = []Todo{}
var myList = MyTodoList{todos}

// créer, modifier, supprimer todo
func (mt *MyTodoList) add(todo Todo) error {
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

// func (mt *MyTodoList) modif(oldText, newText string) error {
// 	for i, t := range mt.todoList {
// 		if t.Text == oldText {
// 			mt.todoList[i].Text = newText
// 			return nil
// 		}
// 	}
// 	return fmt.Errorf("Todo '%s' introuvable", oldText)
// }

func handleClientRequest(w http.ResponseWriter, r *http.Request) {
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

	if !todo.Done {
		err := myList.addTodo(todo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		err := myList.deleteTodo(todo)
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
	list, err := getDb()
	if err != nil {
		log.Fatal(err)
		return
	}
	myList.todoList = list
	json.NewEncoder(w).Encode(myList.todoList)
}

func getDb() ([]Todo, error) {
	var list []Todo

	rows, err := config.Db.Query("SELECT text FROM todos")
	if err != nil {
		return nil, fmt.Errorf("getDb error : %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		todo := Todo{
			Done: false,
		}
		if err := rows.Scan(&todo.Text); err != nil {
			return nil, fmt.Errorf("getDb error : %v", err)
		}
		list = append(list, todo)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getDb error : %v", err)
	}
	return list, nil

}

func (mt *MyTodoList) addTodo(td Todo) error {
	checkpoint := mt.add(td)
	if checkpoint != nil {
		return checkpoint
	}
	result, err := config.Db.Exec("INSERT INTO todos (text) VALUES (?)", td.Text)
	if err != nil {
		return fmt.Errorf("addTodo error : %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("addTodo error : %v", err)
	}
	fmt.Println("id du dernier todo enregistré :", id)
	return nil
}

func (mt *MyTodoList) deleteTodo(td Todo) error {
	checkpoint := mt.delete(td)
	if checkpoint != nil {
		return checkpoint
	}
	result, err := config.Db.Exec("DELETE FROM todos WHERE text LIKE (?)", td.Text)
	if err != nil {
		return fmt.Errorf("deleteTodo error : %v", err)
	}
	fmt.Println(result.RowsAffected())
	return nil
}

func main() {
	db, err := sql.Open("mysql", "adminUser:adminPassword@tcp(db:3306)/todos")
	if err != nil {
		return
	}
	config.Db = db

	fs := http.FileServer(http.Dir("./ihm"))
	http.Handle("/", fs)

	http.HandleFunc("/service", handleClientRequest)
	http.HandleFunc("/todos", getTodos)

	http.ListenAndServe(config.Port, nil)

}
