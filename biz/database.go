package biz

import "webstack/models"

type Database interface {
	AddTodo(td models.Todo) error
	DeleteTodo(td models.Todo) error
	ModifyTodo(td models.Todo) error
	GetTodos() ([]models.Todo, error)
}
