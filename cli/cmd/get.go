package cmd

import (
	"fmt"
	"strconv"
	"webstack/metier/todos"
	"webstack/metier/users"

	"github.com/jedib0t/go-pretty/v6/table"
)

const Chill = 1
const NotChill = 2
const Emergency = 3

const (
	red    = "\x1b[31m"
	yellow = "\x1b[33m"
	green  = "\x1b[32m"
	reset  = "\x1b[0m"
)

type TodoCli struct {
	id       string
	task     string
	priority string
}

func sprintColor(color, text string) string {
	return fmt.Sprintf("%s%s%s", color, text, reset)
}

func NewTodoCli(td todos.Todo) (todocli TodoCli) {
	todocli.id = strconv.Itoa(td.ID)
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
		fmt.Println(err)
		return
	}

	displayed := Todos2TodosCli(list)
	t := table.NewWriter()
	t.SetTitle("My Todolist")
	t.AppendHeader(table.Row{columns[0], columns[1], columns[2]})

	for _, display := range displayed {
		var colorCode string

		switch display.priority {
		case "oula c'est urgent ça":
			colorCode = red
		case "pas particulièrment pressé":
			colorCode = yellow
		case "chill t'as le temps man":
			colorCode = green
		default:
			colorCode = reset
		}

		row := table.Row{
			sprintColor(colorCode, display.id),
			sprintColor(colorCode, display.task),
			sprintColor(colorCode, display.priority),
		}
		t.AppendRow(row)
	}

	fmt.Println(t.Render())
}
