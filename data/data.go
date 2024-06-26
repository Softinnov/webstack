package data

import (
	"database/sql"
	"fmt"

	"webstack/config"
	"webstack/metier/todos"
	"webstack/metier/users"
)

type MysqlServer struct {
}

var db *sql.DB

const ERR_MAILUSED = "email déjà utilisé"
const ERR_NOUSER = "utilisateur introuvable"

func OpenDb(cfg config.Config) (m MysqlServer, err error) {

	urldb := fmt.Sprintf("%s:%s@%s/%s", cfg.Dbusr, cfg.Dbpsw, cfg.Dbsrc, cfg.Db)
	db, err = sql.Open("mysql", urldb)
	if err != nil {
		return m, fmt.Errorf("sql Open() : %v", err)
	}
	return m, nil
}

func CloseDb() error {
	return db.Close()
}

func (m MysqlServer) AddUserDb(user users.User) error {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", users.GetEmail(user)).Scan(&count)
	if err != nil {
		return fmt.Errorf("AddUser error : %v", err)
	}

	if count > 0 {
		return fmt.Errorf(ERR_MAILUSED)
	}
	_, err = db.Exec("INSERT INTO users (email, password) VALUES (?,?)", users.GetEmail(user), users.GetMdp(user))
	if err != nil {
		return fmt.Errorf("AddUser error : %v", err)
	}
	return nil
}

func (m MysqlServer) GetUser(u users.User) (users.User, error) {
	var storedPassword string

	err := db.QueryRow("SELECT password FROM users WHERE email = ?", users.GetEmail(u)).Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return users.User{}, fmt.Errorf(ERR_NOUSER)
		}
		return users.User{}, fmt.Errorf("erreur de connexion à la base de donnée : %v", err)
	}
	u = users.SetMdp(storedPassword)
	return u, nil
}

func (m MysqlServer) GetTodosDb(u users.User) ([]todos.Todo, error) {
	var list []todos.Todo

	rows, err := db.Query("SELECT todoid, text, priority FROM todos JOIN users ON todos.userid = users.userid WHERE users.email = (?)", users.GetEmail(u))
	if err != nil {
		return nil, fmt.Errorf("GetTodos error : %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		todo := todos.Todo{}
		var text string
		if err := rows.Scan(&todo.Id, &text, &todo.Priority); err != nil {
			return nil, fmt.Errorf("GetTodos error : %v", err)
		}
		todo.Task, err = todos.NewTask(text)
		if err != nil {
			return nil, err
		}
		list = append(list, todo)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetTodos error : %v", err)
	}
	return list, nil
}

func (m MysqlServer) AddTodoDb(td todos.Todo, u users.User) error {
	_, err := db.Exec("INSERT INTO todos (text, priority, userid) VALUES (?,?,(SELECT userid FROM users WHERE email = (?)))", todos.GetTask(td.Task), td.Priority, users.GetEmail(u))
	if err != nil {
		return fmt.Errorf("addTodo error : %v", err)
	}
	return nil
}

func (m MysqlServer) DeleteTodoDb(td todos.Todo) error {
	_, err := db.Exec("DELETE FROM todos WHERE todoid = (?)", td.Id)
	if err != nil {
		return fmt.Errorf("deleteTodo error : %v", err)
	}
	return nil
}

func (m MysqlServer) ModifyTodoDb(td todos.Todo) error {
	_, err := db.Exec("UPDATE todos SET text = (?), priority = (?) WHERE todoid = (?)", todos.GetTask(td.Task), td.Priority, td.Id)
	if err != nil {
		return fmt.Errorf("modifyTodo error : %v", err)
	}
	return nil
}
