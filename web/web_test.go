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
	"webstack/metier"
	"webstack/models"
)

func TestEncodejson(t *testing.T) {
	var tests = []struct {
		Id   int
		Text string
	}{
		{Id: 10, Text: "Blabla"},
		{Id: 123, Text: "(/$-_~+)="},
		{Id: 516, Text: ""},
		{Id: 0, Text: "(/$-_]&[~]%)="},
		{Id: 56, Text: "text"},
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
	todos []models.Todo
}

func (f *fakeDb) AddUserDb(u models.User) error {
	return nil
}

func (f *fakeDb) GetUser(u models.User) (models.User, error) {
	return models.User{}, nil
}

func (f *fakeDb) AddTodoDb(td models.Todo, u models.User) error {
	f.todos = append(f.todos, td)
	return nil
}
func (f *fakeDb) DeleteTodoDb(td models.Todo) error {
	for i, t := range f.todos {
		if t.Id == td.Id {
			f.todos = append(f.todos[:i], f.todos[i+1:]...)
			return nil
		}
	}
	return nil

}
func (f *fakeDb) ModifyTodoDb(td models.Todo) error {
	for _, t := range f.todos {
		if t.Id == td.Id {
			t.Text = td.Text
			return nil
		}
	}
	return nil
}
func (f *fakeDb) GetTodosDb(u models.User) (t []models.Todo, e error) {
	t = f.todos
	return t, nil
}

var user models.User

func setupFakeDb() fakeDb {
	db := fakeDb{}

	todo1 := models.Todo{Id: 1, Text: "Faire les courses"}
	todo2 := models.Todo{Id: 2, Text: "Sortir le chien"}

	db.AddTodoDb(todo1, user)
	db.AddTodoDb(todo2, user)

	return db
}

func TestHandleGetTodos(t *testing.T) {
	db := setupFakeDb()
	metier.Init(&db)

	want := db.todos
	req := httptest.NewRequest(http.MethodGet, "/todos", nil)
	w := httptest.NewRecorder()
	HandleGetTodos(w, req)
	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	wantJson, err2 := json.Marshal(want)
	if err2 != nil {
		panic(err2)
	}
	got := string(body)
	if !strings.Contains(got, string(wantJson)) {
		t.Errorf("expected response to contain '%s', but got '%s'", string(wantJson), got)
	}
}

func TestHandleAddTodo(t *testing.T) {
	db := fakeDb{}
	metier.Init(&db)

	var tests = []struct {
		name, entryTxt, entryPrio, want string
	}{
		{"Cas normal", "Blablabla", "2", "Blablabla"},
		{"Chaîne vide", "", "1", "veuillez renseigner du texte"},
		{"Caractères spéciaux autorisés", "(/$-_~+)=", "1", "(/$-_~+)="},
		{"Caractères spéciaux non autorisés", "(/$-_]&[~]%)=", "3", "caractères spéciaux non autorisés"},
		{"Plusieurs espaces en entrée", "    ", "2", "veuillez renseigner du texte"},
		{"Chaîne longue", "Une chaine très longue mais sans caractères spéciaux, d'ailleurs ma mère me dit toujours que je suis spécial, ça va c'est assez long ? Bon aller on va dire que oui", "1", "Une chaine très longue mais sans caractères spéciaux, d'ailleurs ma mère me dit toujours que je suis spécial, ça va c'est assez long ? Bon aller on va dire que oui"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urltxt := fmt.Sprintf("/add?text=%v&priority=%v", url.QueryEscape(tt.entryTxt), tt.entryPrio)
			req := httptest.NewRequest(http.MethodPost, urltxt, nil)
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
	metier.Init(&db)

	var tests = []struct {
		name, entryTxt, entryId, want string
	}{
		{"Cas normal", "Blablabla", "3", "Blablabla"},
		{"Chaîne vide", "", "123", "réessayez ultérieurement"},
		{"Id non numérique", "BlablaASupprimer", "azerty", "erreur de conversion"},
		{"Id vide", "BlablaASupprimer2", "", "erreur de conversion"},
		{"Chaîne longue", "Une chaine très longue mais sans caractères spéciaux, d'ailleurs ma mère me dit toujours que je suis spécial, ça va c'est assez long ? Bon aller on va dire que oui", "10", "Une chaine très longue mais sans caractères spéciaux, d'ailleurs ma mère me dit toujours que je suis spécial, ça va c'est assez long ? Bon aller on va dire que oui"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("/delete?id=%v&text=%v", tt.entryId, url.QueryEscape(tt.entryTxt))
			req := httptest.NewRequest(http.MethodPost, url, nil)
			w := httptest.NewRecorder()
			HandleDeleteTodo(w, req)
			res := w.Result()
			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}
			got := string(body)
			if !strings.Contains(got, tt.want) && !strings.Contains(got, tt.entryId) {
				t.Errorf("expected response to contain '%s' or '%s', but got '%s'", tt.want, tt.entryId, got)
			}
		})
	}
}

func TestHandleModifyTodo(t *testing.T) {
	db := fakeDb{}
	metier.Init(&db)

	var tests = []struct {
		name, entryTxt, entryId, entryPrio, want string
	}{
		{"Cas normal", "Blabliblou", "3", "2", "Blabliblou"},
		{"Chaîne vide", "", "123", "1", "veuillez renseigner du texte"},
		{"Caractères spéciaux autorisés", "(/$-_~+)=", "13", "2", "(/$-_~+)="},
		{"Caractères spéciaux non autorisés", "(/${}-_~+)=", "13", "3", "caractères spéciaux non autorisés"},
		{"Id non numérique", "BlablaAModifier", "azerty", "1", "erreur de conversion"},
		{"Id vide", "BlablaAModifier2", "", "2", "erreur de conversion"},
		{"Plusieurs espaces en entrée", "    ", "56", "2", "veuillez renseigner du texte"},
		{"Chaîne longue", "Une chaine très longue mais sans caractères spéciaux, d'ailleurs ma mère me dit toujours que je suis spécial, ça va c'est assez long ? Bon aller on va dire que oui", "2", "3", "Une chaine très longue mais sans caractères spéciaux, d'ailleurs ma mère me dit toujours que je suis spécial, ça va c'est assez long ? Bon aller on va dire que oui"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("/modify?id=%v&text=%v&priority=%v", tt.entryId, url.QueryEscape(tt.entryTxt), tt.entryPrio)
			req := httptest.NewRequest(http.MethodPost, url, nil)
			w := httptest.NewRecorder()
			HandleModifyTodo(w, req)
			res := w.Result()
			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}
			got := string(body)
			if !strings.Contains(got, tt.want) && !strings.Contains(got, tt.entryId) {
				t.Errorf("expected response to contain '%s' or '%s', but got '%s'", tt.want, tt.entryId, got)
			}
		})
	}
}
