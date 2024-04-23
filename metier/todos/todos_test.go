package todos

import (
	"encoding/json"
	"strings"
	"testing"
	"webstack/metier/users"
)

type fakeDb struct {
	todos []Todo
}

func (f *fakeDb) AddUserDb(u users.User) error {
	return nil
}
func (f *fakeDb) GetUser(u users.User) (users.User, error) {
	return users.User{}, nil
}
func (f *fakeDb) AddTodoDb(td Todo, u users.User) error {
	f.todos = append(f.todos, td)
	return nil
}
func (f *fakeDb) DeleteTodoDb(td Todo) error {
	for i, t := range f.todos {
		if t.ID == td.ID {
			f.todos = append(f.todos[:i], f.todos[i+1:]...)
			return nil
		}
	}

	return nil
}
func (f *fakeDb) ModifyTodoDb(td Todo) error {
	for _, t := range f.todos {
		if t.ID == td.ID {
			t.Task = td.Task
			return nil
		}
	}

	return nil
}
func (f *fakeDb) GetTodosDb(u users.User) (t []Todo, e error) {
	t = f.todos
	return t, nil
}

var user users.User

func setupFakeDb() fakeDb {
	db := fakeDb{}

	task1, _ := NewTask("Faire les courses")
	task2, _ := NewTask("Sortir le chien")
	task3, _ := NewTask("(/$-_~+)=")
	task4, _ := NewTask("Une chaine très longue mais sans caractères spéciaux, d'ailleurs ma mère me dit toujours que je suis spécial, ça va c'est assez long ? Bon aller on va dire que oui")

	todo1 := Todo{ID: 1, Task: task1, Priority: 3}
	todo2 := Todo{ID: 2, Task: task2, Priority: 2}
	todo3 := Todo{ID: 3, Task: task3, Priority: 2}
	todo4 := Todo{ID: 12, Task: task4, Priority: 1}

	db.AddTodoDb(todo1, user)
	db.AddTodoDb(todo2, user)
	db.AddTodoDb(todo3, user)
	db.AddTodoDb(todo4, user)

	return db
}

func TestGetTodos(t *testing.T) {
	db := setupFakeDb()

	err := Init(&db)
	if err != nil {
		t.Fatal(err)
	}

	want := db.todos
	got, err := Get(user)

	if err != nil {
		t.Errorf("Erreur lors de la récupération des données : %v", err)
	}

	if len(got) != len(want) {
		t.Errorf("Expected %d items, but got %d items", len(want), len(got))
		return
	}

	for i := range want {
		if got[i] != want[i] {
			t.Errorf("Expecte item %d to be '%v', but got '%v'", i, want[i], got[i])
		}
	}
}

func TestAddTodo(t *testing.T) {
	db := setupFakeDb()

	err := Init(&db)
	if err != nil {
		t.Fatal(err)
	}

	var tests = []struct {
		name      string
		entryText string
		entryPrio int
		want      string
	}{
		{"Cas normal", "Sortir le chien", 2, db.todos[1].Task.text},
		{"Chaîne vide", "", 1, ErrNoText},
		{"Caractères spéciaux autorisés", "(/$-_~+)=", 1, db.todos[2].Task.text},
		{"Caractères spéciaux non autorisés", "(/$-_]&[~]%)=", 3, ErrSpecialChar},
		{"Plusieurs espaces en entrée", "    ", 2, ErrNoText},
		{"Chaîne longue", "Une chaine très longue mais sans caractères spéciaux, d'ailleurs ma mère me dit toujours que je suis spécial, ça va c'est assez long ? Bon aller on va dire que oui", 1, db.todos[3].Task.text},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := NewTask(tt.entryText)
			got, _ := Add(task, tt.entryPrio, user)
			if err != nil && err.Error() != tt.want {
				t.Errorf("expected response to contain '%s', but got '%s'", tt.want, err)
			}
			if got.Task.text != "" && got.Task.text != tt.want {
				t.Errorf("expected response to contain '%s', but got '%s'", tt.want, got.Task.text)
			}
		})
	}
}

func TestDeleteTodo(t *testing.T) {
	db := fakeDb{}

	err := Init(&db)
	if err != nil {
		t.Fatal(err)
	}

	var tests = []struct {
		name    string
		entryID int
		want    string
	}{
		{"Cas normal", 3, "3"},
		{"Chaîne vide", 123, "123"},
		{"Chaîne longue", 10, "10"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Delete(tt.entryID)
			gotJSON, err2 := json.Marshal(got)
			if err2 != nil {
				panic(err2)
			}
			if (!strings.Contains(string(gotJSON), tt.want)) && !strings.Contains(err.Error(), tt.want) {
				t.Errorf("expected response to contain '%s', but got '%s'", tt.want, err.Error())
			}
		})
	}
}

func TestModifyTodo(t *testing.T) {
	db := setupFakeDb()

	err := Init(&db)
	if err != nil {
		t.Fatal(err)
	}

	var tests = []struct {
		name      string
		entryTxt  string
		entryID   int
		entryPrio int
		want      string
	}{
		{"Cas normal", "Sortir le chien", 3, 2, db.todos[1].Task.text},
		{"Chaîne vide", "", 123, 1, ErrNoText},
		{"Caractères spéciaux autorisés", "(/$-_~+)=", 3, 2, "(/$-_~+)="},
		{"Caractères spéciaux non autorisés", "(/${}-_~+)=", 13, 1, ErrSpecialChar},
		{"Chaîne longue", "Une chaine très longue mais sans caractères spéciaux, d'ailleurs ma mère me dit toujours que je suis spécial, ça va c'est assez long ? Bon aller on va dire que oui", 12, 1, "Une chaine très longue mais sans caractères spéciaux, d'ailleurs ma mère me dit toujours que je suis spécial, ça va c'est assez long ? Bon aller on va dire que oui"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := NewTask(tt.entryTxt)
			got, _ := Modify(task, tt.entryID, tt.entryPrio)
			if err != nil && err.Error() != tt.want {
				t.Errorf("expected response to contain '%s', but got '%s'", tt.want, err)
			}
			if got.Task.text != "" && got.Task.text != tt.want {
				t.Errorf("expected response to contain '%s', but got '%s'", tt.want, got.Task.text)
			}
		})
	}
}
