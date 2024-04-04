package metier

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"webstack/models"
)

var storeTodo DatabaseTodo
var todo models.Todo
var user models.User
var err error

const errSpecialChar = "caractères spéciaux non autorisés : {}[]|"
const errNoText = "veuillez renseigner du texte"
const errNoId = "todo introuvable, réessayez ultérieurement"

func Init(db DatabaseTodo) error {
	if db == nil {
		return fmt.Errorf("db est nil")
	}
	storeTodo = db
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

func sortByPriorityDesc(todos []models.Todo) []models.Todo {
	sort.Slice(todos, func(i, j int) bool {
		return todos[i].Priority > todos[j].Priority
	})
	return todos
}

func GetTodos() ([]models.Todo, error) {
	list, err := storeTodo.GetTodosDb(user)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des données : %v", err)
	}
	listDesc := sortByPriorityDesc(list)
	return listDesc, nil
}

func AddTodo(text string, priorityStr string) (models.Todo, error) {
	if containsSpecialCharacters(text) {
		err = fmt.Errorf(errSpecialChar)
		return models.Todo{}, err
	}
	if containsOnlySpace(text) {
		err = fmt.Errorf(errNoText)
		return models.Todo{}, err
	}
	priority, err := strconv.Atoi(priorityStr)
	if err != nil {
		err = fmt.Errorf("erreur de conversion %v", err)
		return models.Todo{}, err
	}

	todo.Text = text
	todo.Priority = priority
	err = storeTodo.AddTodoDb(todo)
	if err != nil {
		return models.Todo{}, err
	}
	return todo, nil
}

func DeleteTodo(idStr string, text string) (models.Todo, error) {
	if containsOnlySpace(idStr) {
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
	err = storeTodo.DeleteTodoDb(todo)
	if err != nil {
		return models.Todo{}, err
	}
	return todo, nil
}

func ModifyTodo(idStr string, text string, priorityStr string) (models.Todo, error) {
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
	priority, err := strconv.Atoi(priorityStr)
	if err != nil {
		err = fmt.Errorf("erreur de conversion %v", err)
		return models.Todo{}, err
	}

	todo.Id = id
	todo.Text = text
	todo.Priority = priority
	err = storeTodo.ModifyTodoDb(todo)
	if err != nil {
		return models.Todo{}, err
	}
	return todo, nil
}
