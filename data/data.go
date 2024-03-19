package data

import (
	"database/sql"
	"fmt"
)

type Database interface {
	GetDb() ([]Todo, error)
	GetConfig() Config
	AddTodo() error
	DeleteTodo() error
}

type Config struct {
	Port string
	Db   *sql.DB
}

type Todo struct {
	Done bool   `json:"done"`
	Text string `json:"text"`
}

var ServConfig = Config{
	Port: ":5050",
}

func (c Config) GetConfig() Config {
	c.Port = ServConfig.Port
	c.Db = ServConfig.Db
	return c
}

func StartServer() error {
	db, err := sql.Open("mysql", "adminUser:adminPassword@tcp(db:3306)/todos")
	if err != nil {
		return err
	}
	ServConfig.Db = db
	return nil
}

func GetDb() ([]Todo, error) {
	var list []Todo

	rows, err := ServConfig.Db.Query("SELECT text FROM todos")
	if err != nil {
		return nil, fmt.Errorf("getDb error : %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		todo := Todo{
			Done: false,
		}
		if err := rows.Scan(&todo.Text); err != nil {
			return nil, fmt.Errorf("getDb error : %v", err)
		}
		list = append(list, todo)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getDb error : %v", err)
	}
	return list, nil

}

func AddTodo(td Todo) error {
	result, err := ServConfig.Db.Exec("INSERT INTO todos (text) VALUES (?)", td.Text)
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
	result, err := ServConfig.Db.Exec("DELETE FROM todos WHERE text LIKE (?)", td.Text)
	if err != nil {
		return fmt.Errorf("deleteTodo error : %v", err)
	}
	fmt.Println(result.RowsAffected())
	return nil
}
