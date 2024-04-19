package cmd

import (
	"fmt"
	"os"
	"strconv"
	"webstack/metier/todos"
	"webstack/metier/users"
)

const ERR_SUPR = "erreur de suppression de votre tâche"

func Delete(u users.User) {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go delete <id>")
		fmt.Println("")
		fmt.Println("id :             l'identifiant numérique du todo que vous souhaitez supprimer")
		os.Exit(1)
	}
	id, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Arguments invalides :")
		fmt.Println("Usage: go run main.go delete <id>")
		fmt.Println("")
		fmt.Println("id :             l'identifiant numérique du todo que vous souhaitez supprimer")
	}
	_, err = todos.Delete(id)
	if err != nil {
		fmt.Println(ERR_SUPR, err)
		return
	}
	Get(u)
}
