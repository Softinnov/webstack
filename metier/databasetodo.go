package metier

import "webstack/models"

type DatabaseTodo interface {
	AddTodoDb(td models.Todo) error
	DeleteTodoDb(td models.Todo) error
	ModifyTodoDb(td models.Todo) error
	GetTodosDb(u models.User) ([]models.Todo, error)
}
