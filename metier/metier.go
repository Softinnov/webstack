package metier

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"webstack/models"
)

var store Database
var todo models.Todo
var err error

const errSpecialChar = "caractères spéciaux non autorisés : {}[]|"
const errNoText = "veuillez renseigner du texte"
const errNoId = "todo introuvable, réessayez ultérieurement"

func Init(db Database) error {
	if db == nil {
		return fmt.Errorf("db est nil")
	}
	store = db
	return nil
}

func containsSpecialCharacters(s string) bool {
	re := regexp.MustCompile(`[{}[\]|]`)
	return re.MatchString(s)
}

func containsOnlySpace(s string) bool {
	noSpaceText := strings.ReplaceAll(s, " ", "")
	return noSpaceText == ""
}

func GetTodos() ([]models.Todo, error) {
	list, err := store.GetTodosDb()
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des données : %v", err)
	}
	return list, nil
}

func AddTodo(text string) (models.Todo, error) {
	if containsSpecialCharacters(text) {
		err = fmt.Errorf(errSpecialChar)
		return models.Todo{}, err
	}
	if containsOnlySpace(text) {
		err = fmt.Errorf(errNoText)
		return models.Todo{}, err
	}
	todo.Text = text
	err = store.AddTodoDb(todo)
	if err != nil {
		return models.Todo{}, err
	}
	return todo, nil
}

func DeleteTodo(idStr string, text string) (models.Todo, error) {
	if containsOnlySpace(text) || containsOnlySpace(idStr) {
		err = fmt.Errorf(errNoId)
		return models.Todo{}, err
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		err = fmt.Errorf("erreur de conversion %v", err)
		return models.Todo{}, err
	}
	todo.Id = id
	todo.Text = text
	err = store.DeleteTodoDb(todo)
	if err != nil {
		return models.Todo{}, err
	}
	return todo, nil
}

func ModifyTodo(idStr string, text string) (models.Todo, error) {
	if containsOnlySpace(idStr) {
		err = fmt.Errorf(errNoId)
		return models.Todo{}, err
	} else if containsOnlySpace(text) {
		err = fmt.Errorf(errNoText)
		return models.Todo{}, err
	}
	if containsSpecialCharacters(text) {
		err = fmt.Errorf(errSpecialChar)
		return models.Todo{}, err
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		err = fmt.Errorf("erreur de conversion %v", err)
		return models.Todo{}, err
	}
	todo.Id = id
	todo.Text = text
	err = store.ModifyTodoDb(todo)
	if err != nil {
		return models.Todo{}, err
	}
	return todo, nil
}
