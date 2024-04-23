package todos

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"webstack/metier/users"
)

type Task struct {
	text string
}
type Todo struct {
	ID       int  `json:"id"`
	Task     Task `json:"task"`
	Priority int  `json:"priority"`
}

var store DatabaseTodo

const ErrSpecialChar = "caractères spéciaux non autorisés : {}[]|"
const ErrNoText = "veuillez renseigner du texte"
const ErrNoID = "todo introuvable, réessayez ultérieurement"
const ErrGetData = "erreur lors de la récupération des données"
const ErrDBNil = "error db nil"
const ErrPriority = "priorité de tâche invalide"

func Init(db DatabaseTodo) error {
	if db == nil {
		return fmt.Errorf(ErrDBNil)
	}

	store = db

	return nil
}

func NewTask(text string) (t Task, err error) {
	if containsSpecialCharacters(text) {
		err = fmt.Errorf(ErrSpecialChar)
		return t, err
	}

	if containsOnlySpaces(text) {
		err = fmt.Errorf(ErrNoText)
		return t, err
	}

	t.text = text

	return t, nil
}

func GetTask(t Task) string {
	return t.text
}

func containsSpecialCharacters(s string) bool {
	re := regexp.MustCompile(`[{}[\]|]`)
	return re.MatchString(s)
}

func containsOnlySpaces(s string) bool {
	noSpaceText := strings.ReplaceAll(s, " ", "")
	return noSpaceText == ""
}

func sortByPriorityDesc(todos []Todo) []Todo {
	sort.Slice(todos, func(i, j int) bool {
		return todos[i].Priority > todos[j].Priority
	})

	return todos
}

func Get(u users.User) ([]Todo, error) {
	list, err := store.GetTodosDb(u)
	if err != nil {
		return nil, fmt.Errorf("%v : %v", ErrGetData, err)
	}

	listDesc := sortByPriorityDesc(list)

	return listDesc, nil
}

func Add(text Task, priority int, user users.User) (td Todo, err error) {
	if priority < 1 || priority > 3 {
		return td, fmt.Errorf("%v: %v", ErrPriority, err)
	}

	td.Priority = priority
	td.Task = text

	err = store.AddTodoDb(td, user)
	if err != nil {
		return td, err
	}

	return td, nil
}

func Delete(id int) (td Todo, err error) {
	td.ID = id
	err = store.DeleteTodoDb(td)

	if err != nil {
		return td, err
	}

	return td, nil
}

func Modify(text Task, id, priority int) (td Todo, err error) {
	td.ID = id
	td.Task = text
	td.Priority = priority
	err = store.ModifyTodoDb(td)

	if err != nil {
		return td, err
	}

	return td, nil
}
