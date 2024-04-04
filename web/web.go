package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"webstack/metier"
)

func encodejson(w http.ResponseWriter, a any) (any, error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(a)
	if err != nil {
		return nil, fmt.Errorf("erreur d'encodage json : %v", err)
	}
	return a, nil
}

func HandleSignin(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmpassword := r.FormValue("confirmpassword")

	user, err := metier.AddUser(email, password, confirmpassword)
	if err != nil {
		err = fmt.Errorf("erreur : %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = encodejson(w, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Printf("Utilisateur enregistré : %v", user.Email)
	// http.Redirect(w, r, "./todo.html", http.StatusSeeOther)
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := metier.Login(email, password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	_, err = encodejson(w, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Printf("Utilisateur connecté : %v", user.Email)
	// http.Redirect(w, r, "./todo.html", http.StatusSeeOther)
}

func HandleAddTodo(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")
	priority := r.FormValue("priority")

	todo, err := metier.AddTodo(text, priority)
	if err != nil {
		err = fmt.Errorf("erreur ajout de todo : %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = encodejson(w, todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	}
}

func HandleGetTodos(w http.ResponseWriter, r *http.Request) {
	list, err := metier.GetTodos()
	if err != nil {
		http.Error(w, "erreur lors de la récupération des données : réessayez ultérieurement", http.StatusInternalServerError)
		return
	}
	_, err = encodejson(w, list)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
