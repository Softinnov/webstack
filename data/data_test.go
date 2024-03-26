package data

import (
	"testing"
	"webstack/config"
)

// A faire en test d'integration, ne valide pas la logique des fonctions mais l'infrastructure de l'application
func TestOpenDb(t *testing.T) {
	got, _ := OpenDb(config.GetConfig())
	want := MysqlServer{}

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
