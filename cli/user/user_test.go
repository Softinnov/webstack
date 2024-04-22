package user

import (
	"fmt"
	"testing"
	"webstack/metier/todos"
	"webstack/metier/users"
)

type fakeDb struct {
	todos []todos.Todo
	users []users.User
}

func (f *fakeDb) AddUserDb(u users.User) error {
	for _, user := range f.users {
		if users.GetEmail(user) == users.GetEmail(u) {
			return fmt.Errorf("email déjà utilisé")
		}
	}
	f.users = append(f.users, u)
	return nil
}
func (f *fakeDb) GetUser(u users.User) (users.User, error) {
	for _, user := range f.users {
		fmt.Print(user)
		if users.GetEmail(user) == users.GetEmail(u) {
			return user, nil
		}
	}
	return users.User{}, fmt.Errorf("error")
}

func (f *fakeDb) AddTodoDb(td todos.Todo, u users.User) error {
	f.todos = append(f.todos, td)
	return nil
}
func (f *fakeDb) DeleteTodoDb(td todos.Todo) error {
	for i, t := range f.todos {
		if t.Id == td.Id {
			f.todos = append(f.todos[:i], f.todos[i+1:]...)
			return nil
		}
	}
	return nil

}
func (f *fakeDb) ModifyTodoDb(td todos.Todo) error {
	for _, t := range f.todos {
		if t.Id == td.Id {
			t.Task = td.Task
			return nil
		}
	}
	return nil
}
func (f *fakeDb) GetTodosDb(u users.User) (t []todos.Todo, e error) {
	t = f.todos
	return t, nil
}

var user users.User

func setupFakeDb() fakeDb {
	db := fakeDb{}

	task1, _ := todos.NewTask("Faire les courses")
	task2, _ := todos.NewTask("Sortir le chien")
	task3, _ := todos.NewTask("(/$-_~+)=")
	task4, _ := todos.NewTask("Une chaine très longue mais sans caractères spéciaux, d'ailleurs ma mère me dit toujours que je suis spécial, ça va c'est assez long ? Bon aller on va dire que oui")

	mdp, _ := users.HashPassword("123456")
	todo1 := todos.Todo{Id: 1, Task: task1, Priority: 3}
	todo2 := todos.Todo{Id: 2, Task: task2, Priority: 2}
	todo3 := todos.Todo{Id: 3, Task: task3, Priority: 2}
	todo4 := todos.Todo{Id: 12, Task: task4, Priority: 1}
	user, _ = users.NewUser("clem@caramail.fr", mdp)

	db.AddTodoDb(todo1, user)
	db.AddTodoDb(todo2, user)
	db.AddTodoDb(todo3, user)
	db.AddTodoDb(todo4, user)
	db.AddUserDb(user)
	return db
}

func TestLogin(t *testing.T) {
	db := setupFakeDb()
	todos.Init(&db)
	users.Init(&db)

	var tests = []struct {
		name, entryEmail, entryPassword string
		want                            any
	}{
		{"Cas normal", "clem@caramail.fr", "123456", "clem@caramail.fr"},
		{"Email vide", "", "123456", users.ERR_NOMAIL},
		{"Mot de passe incorrect", "clem@caramail.fr", "azerty", users.ERR_BADMDP},
		{"Email invalide", "ma@mail.com", "25mai1995", fmt.Sprint(users.ERR_LOGIN, " : error")},
		{"Mot de passe vide", "clement@caramail.com", "", users.ERR_NOMDP},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewUserCfg(tt.entryEmail, tt.entryPassword)
			err := SaveConfig(u)
			if err != nil {
				t.Errorf("Expected error to be nil but got : %v", err)
			}
			got, gotErr := Login()
			if gotErr != nil {
				fmt.Println(gotErr.Error())
				if gotErr.Error() != tt.want {
					t.Errorf("expected : %v, but got : %v", tt.want, gotErr.Error())
				}
			}
			if gotErr == nil && users.GetEmail(got) != tt.want {
				t.Errorf("expected : %v, but got : %v", tt.want, users.GetEmail(got))
			}

		})
	}
}

func TestSignin(t *testing.T) {
	
}