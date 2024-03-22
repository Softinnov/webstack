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
		return
	}
	defer data.CloseDb()
	biz.Init(msql)
	fs := http.FileServer(http.Dir(cfg.StaticDir))
	http.Handle("/", fs)
	http.HandleFunc("/add", biz.HandleAddRequest)
	http.HandleFunc("/delete", biz.HandleDeleteRequest)
	http.HandleFunc("/modify", biz.HandleModifyRequest)
	http.HandleFunc("/todos", biz.GetTodos)

	http.ListenAndServe(cfg.Port, nil)
}
