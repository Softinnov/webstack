package cmd

import (
	"fmt"
	"os"
	"strconv"
	"webstack/metier/todos"
	"webstack/metier/users"
)

func Modify(u users.User) {
	if len(os.Args) != 5 {
		fmt.Println(`Usage: go run main.go modify <id> "text" <priority>`)
		fmt.Println("")
		fmt.Println("id :             l'identifiant numérique du todo que vous souhaitez modifier")
		fmt.Println("text :           une chaîne de caractère décrivant votre todo")
		fmt.Println("priority :       le niveau de priorité de votre tâche entre 1 et 3, du moins au plus urgent")
		os.Exit(1)
	}
	id, err := strconv.Atoi(os.Args[2])
	text := os.Args[3]
	priority, err2 := strconv.Atoi(os.Args[4])
	if err != nil || err2 != nil {
		fmt.Println("Arguments invalides :")
		fmt.Println(`Usage: go run main.go modify <id> "text" <priority>`)
		fmt.Println("")
		fmt.Println("id :             l'identifiant numérique du todo que vous souhaitez modifier")
		fmt.Println("text :           une chaîne de caractère décrivant votre todo")
		fmt.Println("priority :       le niveau de priorité de votre tâche entre 1 et 3, du moins au plus urgent")
	}
	task, err := todos.NewTask(text)
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = todos.Modify(task, id, priority)
	if err != nil {
		fmt.Println(err.Error())
	}
	Get(u)
}
