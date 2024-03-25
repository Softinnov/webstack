package biz

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"webstack/models"
)

var todo models.Todo
var store Database

func Init(db Database) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}
	store = db
	return nil
}

func encodejson(w http.ResponseWriter, todo models.Todo) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todo)
}

func HandleAddTodo(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")

	if text == "" {
		http.Error(w, "Impossible d'ajouter le todo : veuillez renseigner du texte", http.StatusBadRequest)
		return
	}
	todo.Text = text

	err := store.AddTodo(todo)
	if err != nil {
		http.Error(w, "erreur AddTodo : ", http.StatusInternalServerError)
		return
	}
	encodejson(w, todo)
}

func HandleDeleteTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	text := r.FormValue("text")

	if idStr == "" {
		http.Error(w, "Impossible de supprimer le todo : réessayez ultérieurement", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "erreur de conversion", http.StatusInternalServerError)
		return
	}
	todo.Id = id
	todo.Text = text

	err = store.DeleteTodo(todo)
	if err != nil {
		http.Error(w, "erreur DeleteTodo", http.StatusInternalServerError)
		return
	}
	encodejson(w, todo)
}

func HandleModifyTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	text := r.FormValue("text")

	if idStr == "" || text == "" {
		http.Error(w, "impossible de modifier le todo : réessayez ultérieurement", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "erreur de conversion", http.StatusInternalServerError)
		return
	}
	todo.Id = id
	todo.Text = text

	err = store.ModifyTodo(todo)
	if err != nil {
		http.Error(w, "erreur ModifyTodo", http.StatusInternalServerError)
		return
	}
	encodejson(w, todo)
}

func HandleGetTodos(w http.ResponseWriter, r *http.Request) {
	list, err := store.GetTodos()
	if err != nil {
		http.Error(w, "erreur lors de la récupération des données : réessayez ultérieurement", http.StatusInternalServerError)
		log.Fatal("GetTodos : ", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(list)
}
