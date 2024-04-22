package cmd

import (
	"fmt"
	"os"
	"strconv"
	"webstack/metier/todos"
	"webstack/metier/users"
)

func Add(user users.User) {
	if len(os.Args) != 4 {
		fmt.Println(`Usage: go run main.go add "text" <priority>`)
		fmt.Println("")
		fmt.Println("text :           une chaîne de caractère décrivant votre nouveau todo")
		fmt.Println("priority :       le niveau de priorité de votre tâche entre 1 et 3, du moins au plus urgent")
		os.Exit(1)
	}
	text := os.Args[2]
	priority, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println("Arguments invalides :")
		fmt.Println(`Usage: go run main.go add "text" <priority>`)
		fmt.Println("")
		fmt.Println("text :           une chaîne de caractère décrivant votre nouveau todo")
		fmt.Println("priority :       un nombre entre 1 et 3, du moins au plus urgent")

		os.Exit(1)
	}
	task, err := todos.NewTask(text)
	if err != nil {
		fmt.Print(err.Error())
	}
	_, err = todos.Add(task, priority, user)
	if err != nil {
		fmt.Print(err.Error())
	}
	Get(user)
}
