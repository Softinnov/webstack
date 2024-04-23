package cmd

import (
	"fmt"
	"os"
	"strconv"
	"webstack/metier/todos"
	"webstack/metier/users"
)

const AddArgsLen = 4

func Add(user users.User) {
	if len(os.Args) != AddArgsLen {
		Help("add")
		return
	}

	text := os.Args[2]

	priority, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println(InvArgs)
		Help("add")

		return
	}

	task, err := todos.NewTask(text)
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	_, err = todos.Add(task, priority, user)
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	Get(user)
}
