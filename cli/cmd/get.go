package cmd

import (
	"fmt"
	"webstack/metier/todos"
	"webstack/metier/users"

	"github.com/jedib0t/go-pretty/v6/table"
)

const Chill = 1
const NotChill = 2
const Emergency = 3

type TodoCli struct {
	id       int
	task     string
	priority string
}

func NewTodoCli(td todos.Todo) (todocli TodoCli) {
	todocli.id = td.ID
	todocli.task = todos.GetTask(td.Task)

	switch td.Priority {
	case Chill:
		todocli.priority = "chill t'as le temps man"
	case NotChill:
		todocli.priority = "pas particulièrment pressé"
	case Emergency:
		todocli.priority = "oula c'est urgent ça"
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
	t := table.NewWriter()
	t.SetTitle("My Todolist")
	t.AppendHeader(table.Row{columns[0], columns[1], columns[2]})

	for _, display := range displayed {
		t.AppendRow(table.Row{display.id, display.task, display.priority})
	}

	fmt.Println(t.Render())
}
