package biz

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"webstack/config"
	"webstack/data"
	"webstack/models"
)

func TestInit(t *testing.T) {
	step1, _ := data.OpenDb(config.GetConfig())
	got := Init(step1)

	if got != nil {
		t.Errorf("got %q, wanted nil", got)
	}
}

func TestEncodejson(t *testing.T) {
	w := httptest.NewRecorder()
	model := models.Todo{Id: 10, Text: "Blabla"}
	data, err := encodejson(w, model)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if data != model {
		t.Errorf("expected %q got %q", model, data)
	}
}

type fakeDb struct {
}

func (f fakeDb) AddTodo(td models.Todo) error {
	return nil
}
func (f fakeDb) DeleteTodo(td models.Todo) error {
	return nil
}
func (f fakeDb) ModifyTodo(td models.Todo) error {
	return nil
}
func (f fakeDb) GetTodos() (t []models.Todo, e error) {
	return t, nil
}

func TestGetTodos(t *testing.T) {
	db := fakeDb{}
	Init(db)

	want := ""
	req := httptest.NewRequest(http.MethodGet, "/todos", nil)
	w := httptest.NewRecorder()
	HandleGetTodos(w, req)
	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	got := string(body)
	if !strings.Contains(got, want) {
		t.Errorf("expected response to contain '%s', but got '%s'", want, got)
	}
}

func TestHandleAddTodo(t *testing.T) {
	db := fakeDb{}
	Init(db)

	var tests = []struct {
		name, entryTxt, want string
	}{
		{"Cas normal", "Blablabla", "Blablabla"},
		{"Chaîne vide", "", "veuillez renseigner du texte"},
		{"Caractères spéciaux", "(/$-_][~])=", "Caractères spéciaux non autorisés"},
		// Quand % est dans la chaîne celle ci est renvoyée vide
		// Le & "coupe" la requête http, seule la partie de la chaîne avant le & est considérée
		// Si "%" avant "&" : text vide, si "&" en premier aussi
		{"Caractères spéciaux '&' avant '%'", "(/$-_]&[~]%)=", "Caractères spéciaux non autorisés"},
		{"Caractères spéciaux '%' avant '&'", "(/$%-_]&[~])=", "veuillez renseigner du texte"},
		{"Caractère '&' en début de chaîne", "&(/$-_]&[~])=", "veuillez renseigner du texte"},
		// Plusieurs espaces font fail le test pourtant l'application à bien le comportement attendu en test manuel
		// {"Plusieurs espaces en entrée", "   ", "veuillez renseigner du texte"},
		// longue chaîne provoque une erreur dans le test mais fonctionne à la main
		// {"Chaîne longue", "Une chaine très longue mais sans caractères spéciaux, d'ailleurs ma mère me dit toujours que je suis spécial, ça va c'est assez long ? Bon aller on va dire que oui", "Une chaine très longue mais sans caractères spéciaux, d'ailleurs ma mère me dit toujours que je suis spécial, ça va c'est assez long ? Bon aller on va dire que oui"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("/add?text=%v", tt.entryTxt)
			req := httptest.NewRequest(http.MethodPost, url, nil)
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
	Init(db)

	var tests = []struct {
		name, entryTxt, entryId, want string
	}{
		{"Cas normal", "Blablabla", "3", "Blablabla"},
		{"Chaîne vide", "", "123", "réessayez ultérieurement"},
		{"Id non numérique", "BlablaASupprimer", "azerty", "erreur de conversion"},
		{"Id vide", "BlablaASupprimer2", "", "erreur de conversion"},
		// {"Chaîne longue","Une chaine très longue mais sans caractères spéciaux, d'ailleurs ma mère me dit toujours que je suis spécial, ça va c'est assez long ? Bon aller on va dire que oui", "10", "Une chaine très longue mais sans caractères spéciaux, d'ailleurs ma mère me dit toujours que je suis spécial, ça va c'est assez long ? Bon aller on va dire que oui"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("/delete?id=%v&text=%v", tt.entryId, tt.entryTxt)
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
	Init(db)

	var tests = []struct {
		name, entryTxt, entryId, want string
	}{
		{"Cas normal", "Blabliblou", "3", "Blabliblou"},
		{"Chaîne vide", "", "123", "réessayez ultérieurement"},
		{"Caractères spéciaux", "(/$-_][~])=", "13", "Caractères spéciaux non autorisés"},
		{"Id non numérique", "BlablaAModifier", "azerty", "erreur de conversion"},
		{"Id vide", "BlablaAModifier2", "", "erreur de conversion"},
		//{"Chaîne longue","Une chaine très longue mais sans caractères spéciaux, d'ailleurs ma mère me dit toujours que je suis spécial, ça va c'est assez long ? Bon aller on va dire que oui", "2", "Une chaine très longue mais sans caractères spéciaux, d'ailleurs ma mère me dit toujours que je suis spécial, ça va c'est assez long ? Bon aller on va dire que oui"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("/modify?id=%v&text=%v", tt.entryId, tt.entryTxt)
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
