package cmd

import (
	"fmt"
	"webstack/metier/todos"
	"webstack/metier/users"
)

func Add(text string, priority int, user users.User) {
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
