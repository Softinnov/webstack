package main

import (
	"log"
	"net/http"

	"webstack/config"
	"webstack/data"
	"webstack/metier"
	"webstack/web"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cfg := config.GetConfig()

	msql, err := data.OpenDb(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer data.CloseDb()
	err = metier.Init(msql)
	if err != nil {
		log.Fatal(err)
	}
	fs := http.FileServer(http.Dir(cfg.StaticDir))
	http.Handle("/", fs)
	http.HandleFunc("/add", web.HandleAddTodo)
	http.HandleFunc("/delete", web.HandleDeleteTodo)
	http.HandleFunc("/modify", web.HandleModifyTodo)
	http.HandleFunc("/todos", web.HandleGetTodos)

	http.ListenAndServe(cfg.Port, nil)
}
