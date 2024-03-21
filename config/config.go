package config

import (
	"database/sql"
	"os"
)

type Configuration interface {
	GetConfig()
}
type Config struct {
	Port  string
	Dbsrc string
	Db    *sql.DB
}

var ServConfig = Config{
	Port:  ":5050",
	Dbsrc: os.Getenv("DBS"),
}

func (c Config) GetConfig() Config {
	if ServConfig.Dbsrc == "tcp" {
		c.Dbsrc = "tcp(db:3306)"
	}
	c.Port = ServConfig.Port
	c.Db = ServConfig.Db
	return c
}
