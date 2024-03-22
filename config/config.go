package config

import (
	"os"
)

type Config struct {
	StaticDir string
	Port      string
	Db        string
	Dbsrc     string
	Dbusr     string
	Dbpsw     string
}

var servConfig = Config{
	StaticDir: "./",
	Port:      ":5050",
	Db:        "todos",
	Dbsrc:     "",
	Dbusr:     "adminUser",
	Dbpsw:     "adminPassword",
}

func GetConfig() Config {
	// détermine si base de donnée locale ou dans container
	dbs := os.Getenv("DBS")
	if dbs == "tcp" {
		servConfig.Dbsrc = "tcp(db:3306)"
	}

	// récupére le dossier statique
	dir := os.Getenv("DIR")
	if dir != "" {
		servConfig.StaticDir = dir
	}
	return servConfig
}
