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
		user.Auth(cmd.Get)
	case "add":
		user.Auth(cmd.Add)
	case "delete":
		user.Auth(cmd.Delete)
	case "modify":
		user.Auth(cmd.Modify)
	case "signin":
		user.Signin()
	case "login":
		user.Login()
	case "logout":
		user.ClearUserConfig()
	default:
		fmt.Println("Commande invalide. Utilisez la commande 'help' pour afficher les commandes disponibles.")
		os.Exit(1)
	}
}
