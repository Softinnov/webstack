package metier

import (
	"encoding/json"
	"strings"
	"testing"
	"webstack/config"
	"webstack/data"
	"webstack/models"

	. "github.com/smartystreets/goconvey/convey"
)

func TestInit(t *testing.T) {
	Convey("Test Init(db)", t, func() {
		step1, _ := data.OpenDb(config.GetConfig())
		got := Init(step1)
		Convey("The value should be nil", func() {
			So(got, ShouldBeNil)
		})
	})
	Convey("Test Init(nil)", t, func() {
		got := Init(nil)
		Convey("The value souldn't be nil", func() {
			So(got, ShouldNotBeNil)
		})
	})
}

type fakeDb struct {
	todos []models.Todo
}

func (f *fakeDb) AddTodoDb(td models.Todo) error {
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
func (f *fakeDb) GetTodosDb() (t []models.Todo, e error) {
	t = f.todos
	return t, nil
}

func setupFakeDb() fakeDb {
	db := fakeDb{}

	todo1 := models.Todo{Id: 1, Text: "Faire les courses"}
	todo2 := models.Todo{Id: 2, Text: "Sortir le chien"}

	db.AddTodoDb(todo1)
	db.AddTodoDb(todo2)

	return db
}

func TestGetTodos(t *testing.T) {
	db := setupFakeDb()
	Init(&db)

	want := db.todos
	got, err := GetTodos()

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
	db := fakeDb{}
	Init(&db)

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
			got, err := AddTodo(tt.entryTxt, tt.entryPrio)
			gotJson, err2 := json.Marshal(got)
			if err2 != nil {
				panic(err2)
			}
			if !strings.Contains(string(gotJson), tt.want) && !strings.Contains(err.Error(), tt.want) {
				t.Errorf("expected response to contain '%s', but got '%s'", tt.want, err.Error())
			}
		})
	}
}

func TestDeleteTodo(t *testing.T) {
	db := fakeDb{}
	Init(&db)

	var tests = []struct {
		name, entryTxt, entryId, want string
	}{
		{"Cas normal", "Blablabla", "3", "Blablabla"},
		{"Chaîne vide", "", "123", ""},
		{"Id non numérique", "BlablaASupprimer", "azerty", "erreur de conversion"},
		{"Id vide", "BlablaASupprimer2", "", "réessayez ultérieurement"},
		{"Chaîne longue", "Une chaine très longue mais sans caractères spéciaux, d'ailleurs ma mère me dit toujours que je suis spécial, ça va c'est assez long ? Bon aller on va dire que oui", "10", "Une chaine très longue mais sans caractères spéciaux, d'ailleurs ma mère me dit toujours que je suis spécial, ça va c'est assez long ? Bon aller on va dire que oui"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeleteTodo(tt.entryId, tt.entryTxt)
			gotJson, err2 := json.Marshal(got)
			if err2 != nil {
				panic(err2)
			}
			if (!strings.Contains(string(gotJson), tt.want) || !strings.Contains(string(gotJson), tt.entryId)) && !strings.Contains(err.Error(), tt.want) {
				t.Errorf("expected response to contain '%s', but got '%s'", tt.want, err.Error())
			}
		})
	}
}

func TestModifyTodo(t *testing.T) {
	db := fakeDb{}
	Init(&db)

	var tests = []struct {
		name, entryTxt, entryId, entryPrio, want string
	}{
		{"Cas normal", "Blabliblou", "3", "2", "Blabliblou"},
		{"Chaîne vide", "", "123", "1", "veuillez renseigner du texte"},
		{"Caractères spéciaux autorisés", "(/$-_~+)=", "13", "3", "(/$-_~+)="},
		{"Caractères spéciaux non autorisés", "(/${}-_~+)=", "13", "1", "caractères spéciaux non autorisés"},
		{"Id non numérique", "BlablaAModifier", "azerty", "2", "erreur de conversion"},
		{"Id vide", "BlablaAModifier2", "", "1", "réessayez ultérieurement"},
		{"Plusieurs espaces en entrée", "    ", "56", "2", "veuillez renseigner du texte"},
		{"Chaîne longue", "Une chaine très longue mais sans caractères spéciaux, d'ailleurs ma mère me dit toujours que je suis spécial, ça va c'est assez long ? Bon aller on va dire que oui", "2", "3", "Une chaine très longue mais sans caractères spéciaux, d'ailleurs ma mère me dit toujours que je suis spécial, ça va c'est assez long ? Bon aller on va dire que oui"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ModifyTodo(tt.entryId, tt.entryTxt, tt.entryPrio)
			gotJson, err2 := json.Marshal(got)
			if err2 != nil {
				panic(err2)
			}
			if (!strings.Contains(string(gotJson), tt.want) || !strings.Contains(string(gotJson), tt.entryId)) && !strings.Contains(err.Error(), tt.want) {
				t.Errorf("expected response to contain '%s', but got '%s'", tt.want, err.Error())
			}
		})
	}
}
