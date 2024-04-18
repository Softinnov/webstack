package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"webstack/cli/cmd"
	"webstack/config"
	"webstack/data"
	"webstack/metier/todos"
	"webstack/metier/users"

	_ "github.com/go-sql-driver/mysql"
)

var admin, _ = users.NewUser("clem@calia.com", "123456")

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
		cmd.Get(admin)
	case "add":
		cmd.Add("", 1, admin)
	default:
		fmt.Println("Invalid command. Please use 'help' command to see available commands.")
		os.Exit(1)
	}
}
