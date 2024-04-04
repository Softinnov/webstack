package data

import (
	"database/sql"
	"fmt"

	"webstack/config"
	"webstack/models"

	"golang.org/x/crypto/bcrypt"
)

type MysqlServer struct {
}

var db *sql.DB

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

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

func (m MysqlServer) AddUserDb(u models.User) error {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", u.Email).Scan(&count)
	if err != nil {
		return fmt.Errorf("AddUser error : %v", err)
	}

	if count > 0 {
		return fmt.Errorf("email déjà utilisé")
	}
	u.Mdp, err = hashPassword(u.Mdp)
	if err != nil {
		return fmt.Errorf("hashPassword error : %v", err)
	}
	fmt.Println(u.Mdp)
	_, err = db.Exec("INSERT INTO users (email, password) VALUES (?,?)", u.Email, u.Mdp)
	if err != nil {
		return fmt.Errorf("AddUser error : %v", err)
	}
	return nil
}

func (m MysqlServer) GetUser(u models.User) (models.User, error) {
	var storedPassword string
	var count int

	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", u.Email).Scan(&count)
	if err != nil {
		return models.User{}, fmt.Errorf("count error : %v", err)
	}

	if count == 0 {
		return models.User{}, fmt.Errorf("email introuvable")
	}
	err = db.QueryRow("SELECT password FROM users WHERE email = ?", u.Email).Scan(&storedPassword)
	if err != nil {
		return models.User{}, fmt.Errorf("erreur de connexion à la base de donnée : %v", err)
	}
	if checkPasswordHash(u.Mdp, storedPassword) {
		fmt.Println(u.Mdp, "correspond bien à", storedPassword)
		return u, nil
	} else {
		return models.User{}, fmt.Errorf("mot de passe incorrect")
	}
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
