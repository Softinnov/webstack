package config

import "testing"

func TestGetConfig(t *testing.T) {
	var dbsrc string
	
	var tests = []struct {
		name, dbs, dir string
		wantErr        error
	}{
		{"Config par défaut (dev)", "", "./", nil},
		{"Config prod", "tcp", "./ihm", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("DBS", tt.dbs)
			t.Setenv("DIR", tt.dir)
			if tt.dbs == "tcp" {
				dbsrc = "tcp(db:3306)"
			} else {
				dbsrc = ""
			}
			got, gotErr := GetConfig()
			want := Config{StaticDir: tt.dir, Port: ":5050", Db: "todos", Dbsrc: dbsrc, Dbusr: "adminUser", Dbpsw: "adminPassword"}
			if got != want && gotErr != tt.wantErr {
				t.Errorf("got %q, wanted %q", got, want)
			}
		})
	}
}
