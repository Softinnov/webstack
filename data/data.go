package data

import (
	"database/sql"
	"fmt"

	"webstack/config"
	"webstack/models"
)

type MysqlServer struct {
}

var db *sql.DB

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

func (m MysqlServer) GetTodosDb() ([]models.Todo, error) {
	var list []models.Todo

	rows, err := db.Query("SELECT todoid, text, priority FROM todos")
	if err != nil {
		return nil, fmt.Errorf("GetTodos error : %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		todo := models.Todo{}
		if err := rows.Scan(&todo.Id, &todo.Text, &todo.Priority); err != nil {
			return nil, fmt.Errorf("GetTodos error : %v", err)
		}
		list = append(list, todo)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetTodos error : %v", err)
	}
	return list, nil
}

func (m MysqlServer) AddTodoDb(td models.Todo) error {
	_, err := db.Exec("INSERT INTO todos (text,priority) VALUES (?,?)", td.Text, td.Priority)
	if err != nil {
		return fmt.Errorf("addTodo error : %v", err)
	}
	return nil
}

func (m MysqlServer) DeleteTodoDb(td models.Todo) error {
	_, err := db.Exec("DELETE FROM todos WHERE todoid = (?)", td.Id)
	if err != nil {
		return fmt.Errorf("deleteTodo error : %v", err)
	}
	return nil
}

func (m MysqlServer) ModifyTodoDb(td models.Todo) error {
	_, err := db.Exec("UPDATE todos SET text = (?), priority = (?) WHERE todoid = (?)", td.Text, td.Priority, td.Id)
	if err != nil {
		return fmt.Errorf("modifyTodo error : %v", err)
	}
	return nil
}
