package todos

import (
	"webstack/metier/users"
)

type DatabaseTodo interface {
	AddTodoDb(td Todo, u users.User) error
	DeleteTodoDb(td Todo) error
	ModifyTodoDb(td Todo) error
	GetTodosDb(u users.User) ([]Todo, error)
}
