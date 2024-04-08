package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"webstack/metier"
)

func encodejson(w http.ResponseWriter, a any) (any, error) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(a)
	if err != nil {
		return nil, fmt.Errorf("erreur d'encodage json : %v", err)
	}
	return a, nil
}

func HandleAddTodo(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")
	priority := r.FormValue("priority")

	tokenStr := getActiveCookieTkn(w, r)

	todo, err := metier.AddTodo(text, priority, getUserEmail(tokenStr))
	if err != nil {
		err = fmt.Errorf("erreur ajout de todo : %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = encodejson(w, todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func HandleDeleteTodo(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	text := r.FormValue("text")

	todo, err := metier.DeleteTodo(id, text)
	if err != nil {
		err = fmt.Errorf("erreur suppression de todo : %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = encodejson(w, todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func HandleModifyTodo(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	text := r.FormValue("text")
	priority := r.FormValue("priority")

	todo, err := metier.ModifyTodo(id, text, priority)
	if err != nil {
		err = fmt.Errorf("erreur modification de todo : %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = encodejson(w, todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func HandleGetTodos(w http.ResponseWriter, r *http.Request) {
	tokenStr := getActiveCookieTkn(w, r)

	list, err := metier.GetTodos(getUserEmail(tokenStr))
	if err != nil {
		http.Error(w, "erreur lors de la récupération des données : réessayez ultérieurement", http.StatusInternalServerError)
		return
	}
	_, err = encodejson(w, list)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
