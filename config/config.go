package config

import (
	"os"
)

type Config struct {
	StaticDir string
	Port      string
	Dbsrc     string
}

var servConfig = Config{
	StaticDir: "./",
	Port:      ":5050",
	Dbsrc:     "",
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
