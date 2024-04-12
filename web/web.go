package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"webstack/metier"
)

const ERR_AJOUT = "erreur d'ajout de votre tâche"
const ERR_SUPR = "erreur de suppression de votre tâche"
const ERR_MODIF = "erreur de modification de votre tâche"
const ERR_GETDATA = "erreur lors de la récupération des données"
const ERR_ENCOD = "erreur d'encodage json"

func encodejson(w http.ResponseWriter, a any) (any, error) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(a)
	if err != nil {
		return nil, fmt.Errorf("%v : %v", ERR_ENCOD, err)
	}
	return a, nil
}

func HandleAddTodo(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")
	priority := r.FormValue("priority")

	tokenStr, err := getTokenString(r, COOKIE_NAME)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	todo, err := metier.AddTodo(text, priority, getUserEmail(tokenStr))
	if err != nil {
		err = fmt.Errorf("%v : %v", ERR_AJOUT, err)
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
		err = fmt.Errorf("%v : %v", ERR_SUPR, err)
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
		err = fmt.Errorf("%v : %v", ERR_MODIF, err)
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
	tokenStr, err := getTokenString(r, COOKIE_NAME)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	list, err := metier.GetTodos(getUserEmail(tokenStr))
	if err != nil {
		http.Error(w, ERR_GETDATA, http.StatusInternalServerError)
		return
	}
	_, err = encodejson(w, list)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
