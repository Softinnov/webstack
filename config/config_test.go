package config

import "testing"

func TestGetConfig(t *testing.T) {
	got := GetConfig()
	want := Config{"./", ":5050", "todos", "", "adminUser", "adminPassword"}

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
