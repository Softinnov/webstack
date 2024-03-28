package metier

import "webstack/models"

type Database interface {
	AddTodoDb(td models.Todo) error
	DeleteTodoDb(td models.Todo) error
	ModifyTodoDb(td models.Todo) error
	GetTodosDb() ([]models.Todo, error)
}
