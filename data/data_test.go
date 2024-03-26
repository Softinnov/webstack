package data

import (
	"testing"
	"webstack/config"
)

func TestOpenDb(t *testing.T) {
	got, _ := OpenDb(config.GetConfig())
	want := MysqlServer{}

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
