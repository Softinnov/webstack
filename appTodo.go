package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Todo struct {
	Done string
	Text string
}

func renvoiValeur(w http.ResponseWriter, r *http.Request) {
	// if err := r.ParseForm(); err != nil {
	// 	fmt.Fprintf(w, "ParseForm() err: %v", err)
	// 	return
	// }
	// fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
	done := r.FormValue("check")
	text := r.FormValue("text")

	todo := Todo{
		Done: done,
		Text: text,
	}
	log.Printf("Debug : %v", todo)
	//fmt.Fprintf(w, "Prenom : %s\nNom : %s\n", prenom, nom)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todo)
}

func main() {
	// handle route using handler function
	fs := http.FileServer(http.Dir("./ihm"))
	http.Handle("/", fs)

	http.HandleFunc("/service", renvoiValeur)

	// listen to port
	http.ListenAndServe(":5050", nil)
}
