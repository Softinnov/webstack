package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"webstack/cli/cmd"
	"webstack/cli/user"
	"webstack/config"
	"webstack/data"
	"webstack/metier/todos"
	"webstack/metier/users"

	_ "github.com/go-sql-driver/mysql"
)

const ERRNOTACMD = "Commande invalide. Utilisez la commande 'help' pour afficher les commandes disponibles."

func main() {
	cfg := config.GetConfig()
	msql, err := data.OpenDb(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer data.CloseDb()
	err = todos.Init(msql)
	if err != nil {
		log.Fatal(err)
	}
	err = users.Init(msql)
	if err != nil {
		log.Fatal(err)
	}

	flag.Parse()
	switch flag.Arg(0) {
	case "":
		cmd.Help()
	case "help":
		cmd.Help()
	case "get":
		user.Auth(cmd.Get, user.CFGFILEPATH)
	case "add":
		user.Auth(cmd.Add, user.CFGFILEPATH)
	case "delete":
		user.Auth(cmd.Delete, user.CFGFILEPATH)
	case "modify":
		user.Auth(cmd.Modify, user.CFGFILEPATH)
	case "signin":
		user.Signin(user.CFGFILEPATH)
	case "login":
		user.Login(user.CFGFILEPATH)
	case "logout":
		user.ClearUserConfig(user.CFGFILEPATH)
	default:
		fmt.Println(ERRNOTACMD)
		os.Exit(1)
	}
}
