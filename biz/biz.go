package biz

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"webstack/models"
)

var todo models.Todo
var store Database

func Init(db Database) {
	store = db
}

func HandleAddRequest(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")

	if text == "" {
		return
	}
	todo.Text = text

	err := store.AddTodo(todo)
	if err != nil {
		http.Error(w, "erreur AddTodo : ", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todo)
}

func HandleDeleteRequest(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	text := r.FormValue("text")

	if idStr == "" {
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "erreur de conversion", http.StatusBadRequest)
		return
	}
	todo.Id = id
	todo.Text = text

	err = store.DeleteTodo(todo)
	if err != nil {
		http.Error(w, "erreur DeleteTodo", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todo)
}

func HandleModifyRequest(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	text := r.FormValue("text")

	if idStr == "" || text == "" {
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "erreur de conversion", http.StatusBadRequest)
		return
	}
	todo.Id = id
	todo.Text = text

	err = store.ModifyTodo(todo)
	if err != nil {
		http.Error(w, "erreur ModifyTodo", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todo)
}

func GetTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	list, err := store.GetData()
	if err != nil {
		log.Fatal("getTodos : ", err)
		return
	}
	json.NewEncoder(w).Encode(list)
}
