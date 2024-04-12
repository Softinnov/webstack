package users

type DatabaseUser interface {
	AddUserDb(u User) error
	GetUser(u User) (User, error)
}
