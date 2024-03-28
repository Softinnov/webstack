package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"webstack/metier"
	"webstack/models"
)

var todo models.Todo

func encodejson(w http.ResponseWriter, a any) (any, error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(a)
	if err != nil {
		return nil, fmt.Errorf("erreur d'encodage json : %v", err)
	}
	return a, nil
}

func HandleAddTodo(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")

	_, err := metier.AddTodo(text)
	if err != nil {
		err = fmt.Errorf("erreur ajout de todo : %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	encodejson(w, todo)
}

func HandleDeleteTodo(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	text := r.FormValue("text")

	_, err := metier.DeleteTodo(id, text)
	if err != nil {
		err = fmt.Errorf("erreur suppression de todo : %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	encodejson(w, todo)
}

func HandleModifyTodo(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	text := r.FormValue("text")

	_, err := metier.ModifyTodo(id, text)
	if err != nil {
		err = fmt.Errorf("erreur modification de todo : %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	encodejson(w, todo)
}

func HandleGetTodos(w http.ResponseWriter, r *http.Request) {
	list, err := metier.GetTodos()
	if err != nil {
		http.Error(w, "erreur lors de la récupération des données : réessayez ultérieurement", http.StatusInternalServerError)
		return
	}
	encodejson(w, list)
}
