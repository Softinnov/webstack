package metier

import "webstack/models"

type DatabaseUser interface {
	AddUserDb(u models.User) error
	GetUser(u models.User) (models.User, error)
}
