package main

import (
	"log"
	"net/http"

	"webstack/biz"
	"webstack/config"
	"webstack/data"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cfg := config.GetConfig()

	msql, err := data.OpenDb(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer data.CloseDb()
	err = biz.Init(msql)
	if err != nil {
		log.Fatal(err)
	}
	fs := http.FileServer(http.Dir(cfg.StaticDir))
	http.Handle("/", fs)
	http.HandleFunc("/add", biz.HandleAddTodo)
	http.HandleFunc("/delete", biz.HandleDeleteTodo)
	http.HandleFunc("/modify", biz.HandleModifyTodo)
	http.HandleFunc("/todos", biz.HandleGetTodos)

	http.ListenAndServe(cfg.Port, nil)
}
