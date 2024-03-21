package data

import (
	"database/sql"
	"fmt"

	"webstack/config"
)

type Todo struct {
	Id     int    `json:"id"`
	Text   string `json:"text"`
	Action string `json:"action"`
}

var db *sql.DB

func OpenDb(cfg config.Config) error {
	var err error

	urldb := fmt.Sprintf("adminUser:adminPassword@%s/todos", cfg.Dbsrc)
	fmt.Println(urldb)
	db, err = sql.Open("mysql", urldb)
	if err != nil {
		return fmt.Errorf("sql Open() : %v", err)
	}
	return nil
}

func CloseDb() error {
	return db.Close()
}

func GetTodos() ([]Todo, error) {
	var list []Todo

	rows, err := db.Query("SELECT todoid, text FROM todos")
	if err != nil {
		return nil, fmt.Errorf("GetTodos error : %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		todo := Todo{
			Action: "",
		}
		if err := rows.Scan(&todo.Id, &todo.Text); err != nil {
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
	result, err := db.Exec("INSERT INTO todos (text) VALUES (?)", td.Text)
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
	result, err := db.Exec("DELETE FROM todos WHERE text LIKE (?)", td.Text)
	if err != nil {
		return fmt.Errorf("deleteTodo error : %v", err)
	}
	fmt.Println(result.RowsAffected())
	return nil
}

func ModifyTodo(td Todo) error {
	result, err := db.Exec("UPDATE todos SET text = (?) WHERE todoid = (?)", td.Text, td.Id)
	if err != nil {
		return fmt.Errorf("modifyTodo error : %v", err)
	}
	fmt.Println(result.RowsAffected())
	return nil
}
