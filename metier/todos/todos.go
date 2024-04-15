package todos

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"webstack/metier/users"
)

type Todo struct {
	Id       int    `json:"id"`
	Text     string `json:"text"`
	Priority int    `json:"priority"`
}

var store DatabaseTodo
var todo Todo
var err error
var user users.User

const ERR_SPECIAL_CHAR = "caractères spéciaux non autorisés : {}[]|"
const ERR_NO_TEXT = "veuillez renseigner du texte"
const ERR_NO_ID = "todo introuvable, réessayez ultérieurement"
const ERR_GETDATA = "erreur lors de la récupération des données"
const ERR_CONV = "erreur de conversion"
const ERR_DBNIL = "error db nil"

func Init(db DatabaseTodo) error {
	if db == nil {
		return fmt.Errorf(ERR_DBNIL)
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

func sortByPriorityDesc(todos []Todo) []Todo {
	sort.Slice(todos, func(i, j int) bool {
		return todos[i].Priority > todos[j].Priority
	})
	return todos
}

func Get(email string) ([]Todo, error) {
	user.Email = email
	list, err := store.GetTodosDb(user)
	if err != nil {
		return nil, fmt.Errorf("%v : %v", ERR_GETDATA, err)
	}
	listDesc := sortByPriorityDesc(list)
	return listDesc, nil
}

func Add(text string, priorityStr string, email string) (Todo, error) {
	if containsSpecialCharacters(text) {
		err = fmt.Errorf(ERR_SPECIAL_CHAR)
		return Todo{}, err
	}
	if containsOnlySpace(text) {
		err = fmt.Errorf(ERR_NO_TEXT)
		return Todo{}, err
	}
	priority, err := strconv.Atoi(priorityStr)
	if err != nil {
		err = fmt.Errorf("%v : %v", ERR_CONV, err)
		return Todo{}, err
	}
	user.Email = email
	todo.Text = text
	todo.Priority = priority
	err = store.AddTodoDb(todo, user)
	if err != nil {
		return Todo{}, err
	}
	return todo, nil
}

func Delete(idStr string, text string) (Todo, error) {
	if containsOnlySpace(idStr) {
		err = fmt.Errorf(ERR_NO_ID)
		return Todo{}, err
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		err = fmt.Errorf("%v : %v", ERR_CONV, err)
		return Todo{}, err
	}
	todo.Id = id
	todo.Text = text
	err = store.DeleteTodoDb(todo)
	if err != nil {
		return Todo{}, err
	}
	return todo, nil
}

func Modify(idStr string, text string, priorityStr string) (Todo, error) {
	if containsOnlySpace(idStr) {
		err = fmt.Errorf(ERR_NO_ID)
		return Todo{}, err
	} else if containsOnlySpace(text) {
		err = fmt.Errorf(ERR_NO_TEXT)
		return Todo{}, err
	}
	if containsSpecialCharacters(text) {
		err = fmt.Errorf(ERR_SPECIAL_CHAR)
		return Todo{}, err
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		err = fmt.Errorf("%v : %v", ERR_CONV, err)
		return Todo{}, err
	}
	priority, err := strconv.Atoi(priorityStr)
	if err != nil {
		err = fmt.Errorf("%v : %v", ERR_CONV, err)
		return Todo{}, err
	}

	todo.Id = id
	todo.Text = text
	todo.Priority = priority
	err = store.ModifyTodoDb(todo)
	if err != nil {
		return Todo{}, err
	}
	return todo, nil
}
