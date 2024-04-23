package cmd

import (
	"fmt"
	"os"
	"strconv"
	"webstack/metier/todos"
	"webstack/metier/users"
)

const ModArgsLen = 5
const InvArgs = "Arguments invalides :"

func Modify(u users.User) {
	if len(os.Args) != ModArgsLen {
		Help("modify")
		return
	}

	text := os.Args[3]
	id, err := strconv.Atoi(os.Args[2])
	priority, err2 := strconv.Atoi(os.Args[4])

	if err != nil || err2 != nil {
		fmt.Println(InvArgs)
		Help("modify")

		return
	}

	task, err := todos.NewTask(text)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	_, err = todos.Modify(task, id, priority)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	Get(u)
}
