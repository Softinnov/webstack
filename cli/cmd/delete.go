package cmd

import (
	"fmt"
	"os"
	"strconv"
	"webstack/metier/todos"
	"webstack/metier/users"
)

const ErrSupr = "erreur de suppression de votre tâche"
const DelArgsLen = 3

func Delete(u users.User) {
	if len(os.Args) != DelArgsLen {
		Help("delete")
		return
	}

	id, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println(InvArgs)
		Help("delete")

		return
	}

	_, err = todos.Delete(id)
	if err != nil {
		fmt.Println(ErrSupr, err)
		return
	}

	Get(u)
}
