package main

import (
	"log"
	"net/http"

	"webstack/config"
	"webstack/data"
	"webstack/metier/todos"
	"webstack/metier/users"
	"webstack/web"

	_ "github.com/go-sql-driver/mysql"
)

type FuncHandler struct {
	HandlerFunc func(w http.ResponseWriter, r *http.Request)
}

func (h FuncHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.HandlerFunc(w, r)
}

var userAuth = web.TokenInfo{
	CookieName: web.CookieName,
	PrivateKey: web.SecretKey,
	Auth: web.Auth{
		Name:       "user",
		IsRequired: true,
	},
}

func main() {
	addhandler := FuncHandler{HandlerFunc: web.HandleAddTodo}
	delhandler := FuncHandler{HandlerFunc: web.HandleDeleteTodo}
	modhandler := FuncHandler{HandlerFunc: web.HandleModifyTodo}
	todoshandler := FuncHandler{HandlerFunc: web.HandleGetTodos}

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

	fs := http.FileServer(http.Dir(cfg.StaticDir))
	http.Handle("/", fs)
	http.HandleFunc("/add", web.WrapAuth(addhandler, userAuth))
	http.HandleFunc("/delete", web.WrapAuth(delhandler, userAuth))
	http.HandleFunc("/modify", web.WrapAuth(modhandler, userAuth))
	http.HandleFunc("/todos", web.WrapAuth(todoshandler, userAuth))
	http.HandleFunc("/signin", web.HandleSignin)
	http.HandleFunc("/login", web.HandleLogin)
	http.HandleFunc("/logout", web.HandleLogout)

	err = http.ListenAndServe(cfg.Port, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer data.CloseDB()
}
