package cmd

import (
	"fmt"
	"strings"
	"webstack/metier/todos"
	"webstack/metier/users"
	"webstack/web"
)

func Get(u users.User) {
	columns := []string{"ID", "Titre", "Priorité"}
	list, err := todos.Get(u)
	if err != nil {
		err = fmt.Errorf("error get : %v", err)
		fmt.Println(err)
	}
	displayed := web.Todos2TodosView(list)
	fmt.Printf("%-10s %-20s %-10s\n", columns[0], columns[1], columns[2])
	fmt.Println(strings.Repeat("-", 40))
	for _, todo := range displayed {
		fmt.Printf("%-10d %-20s %-10d\n", todo.Id, todo.Task, todo.Priority)
	}
}
