package user

import (
	"fmt"
	"testing"
	"webstack/metier/todos"
	"webstack/metier/users"
)

const TestCfg = "../.cfg/config.json"

type fakeDB struct {
	todos []todos.Todo
	users []users.User
}

func (f *fakeDB) AddUserDb(u users.User) error {
	for _, user := range f.users {
		if users.GetEmail(user) == users.GetEmail(u) {
			return fmt.Errorf("email déjà utilisé")
		}
	}

	f.users = append(f.users, u)

	return nil
}
func (f *fakeDB) GetUser(u users.User) (users.User, error) {
	for _, user := range f.users {
		fmt.Print(user)

		if users.GetEmail(user) == users.GetEmail(u) {
			return user, nil
		}
	}

	return users.User{}, fmt.Errorf("error")
}

func (f *fakeDB) AddTodoDb(td todos.Todo, u users.User) error {
	f.todos = append(f.todos, td)
	return nil
}
func (f *fakeDB) DeleteTodoDb(td todos.Todo) error {
	for i, t := range f.todos {
		if t.ID == td.ID {
			f.todos = append(f.todos[:i], f.todos[i+1:]...)
			return nil
		}
	}

	return nil
}
func (f *fakeDB) ModifyTodoDb(td todos.Todo) error {
	for _, t := range f.todos {
		if t.ID == td.ID {
			t.Task = td.Task
			return nil
		}
	}

	return nil
}
func (f *fakeDB) GetTodosDb(u users.User) (t []todos.Todo, e error) {
	t = f.todos
	return t, nil
}

var user users.User

func setupFakeDB() fakeDB {
	db := fakeDB{}

	task1, _ := todos.NewTask("Faire les courses")
	task2, _ := todos.NewTask("Sortir le chien")
	task3, _ := todos.NewTask("(/$-_~+)=")
	task4, _ := todos.NewTask("Une chaine très longue mais sans caractères spéciaux, d'ailleurs ma mère me dit toujours que je suis spécial, ça va c'est assez long ? Bon aller on va dire que oui")

	mdp, _ := users.HashPassword("123456")
	todo1 := todos.Todo{ID: 1, Task: task1, Priority: 3}
	todo2 := todos.Todo{ID: 2, Task: task2, Priority: 2}
	todo3 := todos.Todo{ID: 3, Task: task3, Priority: 2}
	todo4 := todos.Todo{ID: 12, Task: task4, Priority: 1}
	user, _ = users.NewUser("clem@caramail.fr", mdp)

	db.AddTodoDb(todo1, user)
	db.AddTodoDb(todo2, user)
	db.AddTodoDb(todo3, user)
	db.AddTodoDb(todo4, user)
	db.AddUserDb(user)

	return db
}

// Sans doute plus interessant de tester SaveConfig, LoadConfig, ClearUserConfig et NewUserCfg que Login et Signin
func TestUserConfig(t *testing.T) {
	var tests = []struct {
		name, entryEmail, entryPassword string
	}{
		{"Cas normal", "mail2018@mail.com", "29mai1995"},
		{"Mots de passes différents", "mail2019@mail.com", "29mai1995"},
		{"Email vide", "", "12azerty"},
		{"Email invalide", "mail2018mailcom", "29mai1995"},
		{"Mot de passe trop court", "mail@mail.com", "azey"},
		{"Email déjà utilisé", "mail20@mail.com", "2mai1995"},
		{"Mot de passe vide", "clem@caramail.com", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewUserCfg(tt.entryEmail, tt.entryPassword)

			err := SaveConfig(u, TestCfg)
			if err != nil {
				t.Errorf("Expected error to be nil but got : %v", err)
			}

			configData, err := LoadConfig(TestCfg)
			if err != nil {
				t.Errorf("expected error to be nil but got : %v", err)
			}
			if u != configData {
				t.Errorf("expected : %v, but got : %v", u, configData)
			}
		})
	}
}

func TestClearUserConfig(t *testing.T) {
	db := setupFakeDB()
	err := todos.Init(&db)
	err2 := users.Init(&db)

	if err != nil || err2 != nil {
		t.Fatalf("Expected error to be nil but got : %v / %v", err, err2)
	}

	var tests = []struct {
		name, entryEmail, entryPassword string
	}{
		{"Cas normal", "mail2018@mail.com", "29mai1995"},
		{"Mots de passes différents", "mail2019@mail.com", "29mai1995"},
		{"Email vide", "", "12azerty"},
		{"Email invalide", "mail2018mailcom", "29mai1995"},
		{"Email déjà utilisé", "mail20@mail.com", "2mai1995"},
		{"Mot de passe vide", "clem@caramail.com", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewUserCfg(tt.entryEmail, tt.entryPassword)

			err := SaveConfig(u, TestCfg)
			if err != nil {
				t.Errorf("Expected error to be nil but got : %v", err)
			}

			configData, err := LoadConfig(TestCfg)
			if err != nil {
				t.Errorf("expected error to be nil but got : %v", err)
			} else if configData.Email == "" && configData.Mdp == "" {
				t.Errorf("expected : %v, but got nothing", configData)
			}

			err = ClearUserConfig(TestCfg)
			if err != nil {
				t.Errorf("Expected error to be nil but got : %v", err)
			}

			emptyConfigData, err := LoadConfig(TestCfg)
			if (emptyConfigData.Email != "" && emptyConfigData.Mdp != "") || err != nil {
				t.Errorf("expected : %v to be empty but it's not", emptyConfigData)
			}
		})
	}
}
