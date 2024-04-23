package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
	"webstack/metier/todos"
	"webstack/metier/users"
)

func TestEncodejson(t *testing.T) {
	var tests = []struct {
		ID   int
		Text string
	}{
		{ID: 10, Text: "Blabla"},
		{ID: 123, Text: "(/$-_~+)="},
		{ID: 516, Text: ""},
		{ID: 0, Text: "(/$-_]&[~]%)="},
		{ID: 56, Text: "text"},
	}

	for i, tt := range tests {
		t.Run(tt.Text, func(t *testing.T) {
			w := httptest.NewRecorder()
			model := tests[i]
			data, err := encodejson(w, model)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}
			if data != model {
				t.Errorf("expected %q got %q", model, data)
			}
		})
	}
}

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
		if t.ID == td.ID {
			f.todos = append(f.todos[:i], f.todos[i+1:]...)
			return nil
		}
	}

	return nil
}
func (f *fakeDb) ModifyTodoDb(td todos.Todo) error {
	for _, t := range f.todos {
		if t.ID == td.ID {
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

func TestHandleSignin(t *testing.T) {
	db := setupFakeDb()

	err := todos.Init(&db)
	err2 := users.Init(&db)

	if err != nil || err2 != nil {
		t.Fatalf("Expected error to be nil but got : %v / %v", err, err2)
	}

	var tests = []struct {
		name, entryEmail, entryPassword, confirmPassword, want string
	}{
		{"Cas normal", "clem@mail.fr", "123456", "123456", ""},
		{"Email vide", "", "123456", "123456", "l'email ne peut pas être vide"},
		{"Mots de passe différents", "clem@caramail.fr", "azerty", "azery", "mots de passe différents"},
		{"Mot de passe trop court", "clem@mail.com", "123", "123", "mot de passe trop court (6 caractères minimum)"},
		{"Email invalide", "mamail.com", "25mai1995", "25mai1995", "email invalide"},
		{"Mot de passe vide", "clement@caramail.com", "", "", "le mot de passe ne peut pas être vide"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urltxt := fmt.Sprintf("/signin?email=%v&password=%v&confirmpassword=%v", url.QueryEscape(tt.entryEmail), tt.entryPassword, tt.confirmPassword)
			req := httptest.NewRequest(http.MethodPost, urltxt, http.NoBody)
			w := httptest.NewRecorder()
			HandleSignin(w, req)
			res := w.Result()
			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}
			got := string(body)
			if !strings.Contains(got, tt.want) {
				t.Errorf("expected response to contain '%s', but got '%s'", tt.want, got)
			}
		})
	}
}
func TestHandleLogin(t *testing.T) {
	db := setupFakeDb()

	err := todos.Init(&db)
	err2 := users.Init(&db)

	if err != nil || err2 != nil {
		t.Fatalf("Expected error to be nil but got : %v / %v", err, err2)
	}

	var tests = []struct {
		name, entryEmail, entryPassword, want string
	}{
		{"Cas normal", "clem@caramail.fr", "123456", ""},
		{"Email vide", "", "123456", "l'email ne peut pas être vide"},
		{"Mot de passe incorrect", "clem@caramail.fr", "azerty", "mot de passe incorrect"},
		{"Email invalide", "ma@mail.com", "25mai1995", "échec du login"},
		{"Mot de passe vide", "clement@caramail.com", "", "le mot de passe ne peut pas être vide"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urltxt := fmt.Sprintf("/login?email=%v&password=%v", url.QueryEscape(tt.entryEmail), tt.entryPassword)
			req := httptest.NewRequest(http.MethodPost, urltxt, http.NoBody)
			w := httptest.NewRecorder()
			HandleLogin(w, req)
			res := w.Result()
			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}
			got := string(body)
			if !strings.Contains(got, tt.want) {
				t.Errorf("expected response to contain '%s', but got '%s'", tt.want, got)
			}
		})
	}
}

func TestHandleLogout(t *testing.T) {
	db := setupFakeDb()

	err := users.Init(&db)

	if err != nil {
		t.Fatalf("Expected error to be nil but got : %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/logout", http.NoBody)
	req.AddCookie(&http.Cookie{
		Name:  CookieName,
		Value: "token",
	})

	w := httptest.NewRecorder()
	HandleLogout(w, req)
	res := w.Result()
	got := res.Cookies()

	defer res.Body.Close()

	want := &http.Cookie{
		Name:    CookieName,
		Value:   "",
		MaxAge:  -1,
		Expires: time.Now().Add(-time.Hour),
		Path:    "/",
	}
	found := false

	for _, cookie := range got {
		if cookie.Name == want.Name {
			if cookie.Value == want.Value {
				found = true
				break
			}
		}
	}

	if !found {
		t.Errorf("expected response to contain '%s', but got '%s'", want, got)
	}
}

func TestHandleGetTodos(t *testing.T) {
	db := setupFakeDb()

	err := todos.Init(&db)
	err2 := users.Init(&db)

	if err != nil || err2 != nil {
		t.Fatalf("Expected error to be nil but got : %v / %v", err, err2)
	}

	token := jsonwebToken(user)

	want := Todos2TodosView(db.todos)
	req := httptest.NewRequest(http.MethodGet, "/todos", http.NoBody)
	req.AddCookie(&http.Cookie{
		Name:  CookieName,
		Value: token,
	})

	w := httptest.NewRecorder()
	HandleGetTodos(w, req)
	res := w.Result()

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	wantJSON, err2 := json.Marshal(want)
	if err2 != nil {
		panic(err2)
	}

	got := string(body)
	if !strings.Contains(got, string(wantJSON)) {
		t.Errorf("expected response to contain '%s', but got '%s'", string(wantJSON), got)
	}
}

func TestHandleAddTodo(t *testing.T) {
	db := setupFakeDb()

	err := todos.Init(&db)
	err2 := users.Init(&db)

	if err != nil || err2 != nil {
		t.Fatalf("Expected error to be nil but got : %v / %v", err, err2)
	}

	token := jsonwebToken(user)

	var tests = []struct {
		name, entryTxt, entryPrio, want string
	}{
		{"Cas normal", "Sortir le chien", "2", todos.GetTask(db.todos[1].Task)},
		{"Chaîne vide", "", "1", "veuillez renseigner du texte"},
		{"Caractères spéciaux autorisés", "(/$-_~+)=", "2", todos.GetTask(db.todos[2].Task)},
		{"Caractères spéciaux non autorisés", "(/$-_]&[~]%)=", "3", "caractères spéciaux non autorisés"},
		{"Plusieurs espaces en entrée", "    ", "2", "veuillez renseigner du texte"},
		{"Chaîne longue", "Une chaine très longue mais sans caractères spéciaux, d'ailleurs ma mère me dit toujours que je suis spécial, ça va c'est assez long ? Bon aller on va dire que oui", "1", todos.GetTask(db.todos[3].Task)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urltxt := fmt.Sprintf("/add?task=%v&priority=%v", url.QueryEscape(tt.entryTxt), tt.entryPrio)
			req := httptest.NewRequest(http.MethodPost, urltxt, http.NoBody)
			req.AddCookie(&http.Cookie{
				Name:  CookieName,
				Value: token,
			})
			w := httptest.NewRecorder()
			HandleAddTodo(w, req)
			res := w.Result()
			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}
			got := string(body)
			if !strings.Contains(got, tt.want) {
				t.Errorf("expected response to contain '%s', but got '%s'", tt.want, got)
			}
		})
	}
}

func TestHandleDeleteTodo(t *testing.T) {
	db := fakeDb{}

	err := todos.Init(&db)
	err2 := users.Init(&db)

	if err != nil || err2 != nil {
		t.Fatalf("Expected error to be nil but got : %v / %v", err, err2)
	}

	token := jsonwebToken(user)

	var tests = []struct {
		name, entryTxt, entryID, want string
	}{
		{"Cas normal", "Blablabla", "3", "Blablabla"},
		{"Chaîne vide", "", "123", "réessayez ultérieurement"},
		{"Id non numérique", "BlablaASupprimer", "azerty", ErrConv},
		{"Id vide", "BlablaASupprimer2", "", ErrConv},
		{"Chaîne longue", "Une chaine très longue mais sans caractères spéciaux, d'ailleurs ma mère me dit toujours que je suis spécial, ça va c'est assez long ? Bon aller on va dire que oui", "10", "Une chaine très longue mais sans caractères spéciaux, d'ailleurs ma mère me dit toujours que je suis spécial, ça va c'est assez long ? Bon aller on va dire que oui"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urltxt := fmt.Sprintf("/delete?id=%v&text=%v", tt.entryID, url.QueryEscape(tt.entryTxt))
			fmt.Print(urltxt)
			req := httptest.NewRequest(http.MethodPost, urltxt, http.NoBody)
			req.AddCookie(&http.Cookie{
				Name:  CookieName,
				Value: token,
			})
			w := httptest.NewRecorder()
			HandleDeleteTodo(w, req)
			res := w.Result()
			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			fmt.Print(body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}
			got := string(body)
			fmt.Print(got)
			if !strings.Contains(got, tt.want) && !strings.Contains(got, tt.entryID) {
				t.Errorf("expected response to contain '%s' or '%s', but got '%s'", tt.want, tt.entryID, got)
			}
		})
	}
}

func TestHandleModifyTodo(t *testing.T) {
	db := setupFakeDb()

	err := todos.Init(&db)
	err2 := users.Init(&db)

	if err != nil || err2 != nil {
		t.Fatalf("Expected error to be nil but got : %v / %v", err, err2)
	}

	token := jsonwebToken(user)

	var tests = []struct {
		name, entryTxt, entryID, entryPrio, want string
	}{
		{"Cas normal", "Sortir le chien", "2", "2", todos.GetTask(db.todos[1].Task)},
		{"Caractères spéciaux autorisés", "(/$-_~+)=", "3", "2", todos.GetTask(db.todos[2].Task)},
		{"Chaîne vide", "", "123", "1", "veuillez renseigner du texte"},
		{"Caractères spéciaux autorisés", "(/$-_~+)=", "13", "2", todos.GetTask(db.todos[2].Task)},
		{"Caractères spéciaux non autorisés", "(/${}-_~+)=", "13", "3", "caractères spéciaux non autorisés"},
		{"Id non numérique", "BlablaAModifier", "azerty", "1", ErrConv},
		{"Id vide", "BlablaAModifier2", "", "2", ErrConv},
		{"Plusieurs espaces en entrée", "    ", "56", "2", "veuillez renseigner du texte"},
		{"Chaîne longue", "Une chaine très longue mais sans caractères spéciaux, d'ailleurs ma mère me dit toujours que je suis spécial, ça va c'est assez long ? Bon aller on va dire que oui", "12", "1", todos.GetTask(db.todos[3].Task)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urltxt := fmt.Sprintf("/modify?id=%v&task=%v&priority=%v", tt.entryID, url.QueryEscape(tt.entryTxt), tt.entryPrio)
			req := httptest.NewRequest(http.MethodPost, urltxt, http.NoBody)
			req.AddCookie(&http.Cookie{
				Name:  CookieName,
				Value: token,
			})
			w := httptest.NewRecorder()
			HandleModifyTodo(w, req)
			res := w.Result()
			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}
			got := string(body)
			if !strings.Contains(got, tt.want) {
				t.Errorf("expected response to contain '%s', but got '%s'", tt.want, got)
			}
		})
	}
}
