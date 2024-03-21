package data

import (
	"database/sql"
	"fmt"

	"webstack/config"
)

type Database interface {
	GetTodos() ([]Todo, error)
	AddTodo() error
	DeleteTodo() error
}

type Todo struct {
	Done bool   `json:"done"`
	Text string `json:"text"`
}

func OpenDb() (*sql.DB, error) {
	// config.ServConfig.Dbsrc = os.Getenv("DBS")
	// if config.ServConfig.Dbsrc == "tcp" {
	// 	config.ServConfig.Dbsrc = "tcp(db:3306)"
	// }
	dbsrc := config.ServConfig.Dbsrc
	urldb := fmt.Sprintf("adminUser:adminPassword@%s/todos", dbsrc)
	fmt.Println(urldb)
	db, err := sql.Open("mysql", urldb)
	if err != nil {
		return nil, fmt.Errorf("sql Open() : %v", err)
	}
	return db, nil
}

func GetTodos() ([]Todo, error) {
	var list []Todo

	rows, err := config.ServConfig.Db.Query("SELECT text FROM todos")
	if err != nil {
		return nil, fmt.Errorf("GetTodos error : %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		todo := Todo{
			Done: false,
		}
		if err := rows.Scan(&todo.Text); err != nil {
			return nil, fmt.Errorf("GetTodos error : %v", err)
		}
		list = append(list, todo)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetTodos error : %v", err)
	}
	return list, nil

}

func AddTodo(td Todo) error {
	result, err := config.ServConfig.Db.Exec("INSERT INTO todos (text) VALUES (?)", td.Text)
	if err != nil {
		return fmt.Errorf("addTodo error : %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("addTodo error : %v", err)
	}
	fmt.Println("id du dernier todo enregistré :", id)
	return nil
}

func DeleteTodo(td Todo) error {
	result, err := config.ServConfig.Db.Exec("DELETE FROM todos WHERE text LIKE (?)", td.Text)
	if err != nil {
		return fmt.Errorf("deleteTodo error : %v", err)
	}
	fmt.Println(result.RowsAffected())
	return nil
}
