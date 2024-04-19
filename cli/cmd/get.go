package cmd

import (
	"fmt"
	"strings"
	"webstack/metier/todos"
	"webstack/metier/users"
)

type TodoCli struct {
	id       int
	task     string
	priority string
}

func NewTodoCli(td todos.Todo) (todocli TodoCli) {
	todocli.id = td.Id
	todocli.task = todos.GetTask(td.Task)
	if td.Priority == 3 {
		todocli.priority = "urgent"
	} else if td.Priority == 2 {
		todocli.priority = "pas particulièrment pressé"
	} else if td.Priority == 1 {
		todocli.priority = "chill man"
	}
	return todocli
}

func Todos2TodosCli(list []todos.Todo) (displayedList []TodoCli) {
	for _, todo := range list {
		displayedList = append(displayedList, NewTodoCli(todo))
	}
	return displayedList
}

func Get(u users.User) {
	columns := []string{"Id", "Todo", "Priorité"}
	list, err := todos.Get(u)
	if err != nil {
		err = fmt.Errorf("error get : %v", err)
		fmt.Println(err)
		return
	}
	displayed := Todos2TodosCli(list)
	fmt.Printf("%-10s %-20s %-10s\n", columns[0], columns[1], columns[2])
	fmt.Println(strings.Repeat("-", 60))
	for _, todo := range displayed {
		fmt.Printf("%-10d %-20s %-10s\n", todo.id, todo.task, todo.priority)
	}
}
