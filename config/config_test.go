package config

import "testing"

func TestGetConfig(t *testing.T) {
	got := GetConfig()
	want := Config{StaticDir: "./", Port: ":5050", Db: "todos", Dbsrc: "", Dbusr: "adminUser", Dbpsw: "adminPassword"}

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
