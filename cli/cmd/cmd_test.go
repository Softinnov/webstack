package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
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

func TestGet(t *testing.T) {
	db := setupFakeDb()
	todos.Init(&db)
	users.Init(&db)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	defer func() {
		os.Stdout = old
		w.Close()
	}()
	Get(db.users[0])

	outCh := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outCh <- buf.String()
	}()
	w.Close()

	fmt.Print(old)
	want := todos.GetTask(db.todos[0].Task)
	want2 := todos.GetTask(db.todos[1].Task)
	actual := <-outCh

	if !strings.Contains(actual, want) || !strings.Contains(actual, want2) {
		t.Errorf("expected : %s and %s, but got : %s", want, want2, actual)
	}
}

func TestAdd(t *testing.T) {
	db := setupFakeDb()
	todos.Init(&db)
	users.Init(&db)

	var tests = []struct {
		name      string
		entryText string
		entryPrio string
		want      string
	}{
		{"Cas normal", "Sortir le chien", "2", todos.GetTask(db.todos[1].Task)},
		{"Chaîne vide", "", "1", todos.ERR_NO_TEXT},
		{"Caractères spéciaux autorisés", "(/$-_~+)=", "1", todos.GetTask(db.todos[2].Task)},
		{"Caractères spéciaux non autorisés", "(/$-_]&[~]%)=", "3", todos.ERR_SPECIAL_CHAR},
		{"Plusieurs espaces en entrée", "    ", "2", todos.ERR_NO_TEXT},
		{"Chaîne longue", "Une chaine très longue mais sans caractères spéciaux, d'ailleurs ma mère me dit toujours que je suis spécial, ça va c'est assez long ? Bon aller on va dire que oui", "1", todos.GetTask(db.todos[3].Task)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = []string{"go run main.go", "add", tt.entryText, tt.entryPrio}
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			defer func() {
				os.Stdout = old
				w.Close()
			}()
			Add(db.users[0])
			outCh := make(chan string)
			go func() {
				var buf bytes.Buffer
				io.Copy(&buf, r)
				outCh <- buf.String()
			}()
			w.Close()

			actual := <-outCh

			if !strings.Contains(actual, tt.want) {
				t.Errorf("expected : %s, but got : %s", tt.want, actual)
			}

		})
	}
}

func TestDelete(t *testing.T) {
	db := setupFakeDb()
	todos.Init(&db)
	users.Init(&db)

	os.Args = []string{"go run main.go", "delete", "2"}
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	defer func() {
		os.Stdout = old
		w.Close()
	}()

	Delete(db.users[0])

	outCh := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outCh <- buf.String()
	}()
	w.Close()

	dontWant := "Sortir le chien"
	actual := <-outCh

	if strings.Contains(actual, dontWant) {
		t.Errorf("expected : %s, but got : %s", dontWant, actual)
	}
}

func TestModify(t *testing.T) {
	db := setupFakeDb()
	todos.Init(&db)
	users.Init(&db)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	defer func() {
		os.Stdout = old
		w.Close()
	}()

	os.Args = []string{"go run main.go", "modify", "2", "Todo modifié", "1"}
	fmt.Println(os.Args)
	Modify(db.users[0])

	outCh := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outCh <- buf.String()
		fmt.Println(outCh)
	}()
	w.Close()

	want := "Todo modifié"
	actual := <-outCh

	if !strings.Contains(actual, want) {
		t.Errorf("expected : %s, but got : %s", want, actual)
	}
}
