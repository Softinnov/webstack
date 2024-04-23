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

const ErrNotACmd = "Commande invalide. Utilisez la commande 'help' pour afficher les commandes disponibles."

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	msql, err := data.OpenDB(&cfg)
	if err != nil {
		log.Fatal(err)
	}

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
		cmd.Help("help")
	case "help":
		cmd.Help("help")
	case "get":
		user.Auth(cmd.Get, user.CfgFilePath)
	case "add":
		user.Auth(cmd.Add, user.CfgFilePath)
	case "delete":
		user.Auth(cmd.Delete, user.CfgFilePath)
	case "modify":
		user.Auth(cmd.Modify, user.CfgFilePath)
	case "signin":
		_, err = user.Signin(user.CfgFilePath)
		if err != nil {
			fmt.Println(err)
			return
		}
	case "login":
		_, err = user.Login(user.CfgFilePath)
		if err != nil {
			fmt.Println(err)
			return
		}
	case "logout":
		err = user.ClearUserConfig(user.CfgFilePath)
		if err != nil {
			fmt.Println(err)
			return
		}
	default:
		fmt.Println(ErrNotACmd)
		os.Exit(1)
	}
	defer data.CloseDB()
}
