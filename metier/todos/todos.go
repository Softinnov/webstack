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
	Id       int  `json:"id"`
	Task     Task `json:"task"`
	Priority int  `json:"priority"`
}

var store DatabaseTodo

const ERR_SPECIAL_CHAR = "caractères spéciaux non autorisés : {}[]|"
const ERR_NO_TEXT = "veuillez renseigner du texte"
const ERR_NO_ID = "todo introuvable, réessayez ultérieurement"
const ERR_GETDATA = "erreur lors de la récupération des données"
const ERR_DBNIL = "error db nil"

func Init(db DatabaseTodo) error {
	if db == nil {
		return fmt.Errorf(ERR_DBNIL)
	}
	store = db
	return nil
}

func NewTask(text string) (t Task, err error) {
	if containsSpecialCharacters(text) {
		err = fmt.Errorf(ERR_SPECIAL_CHAR)
		return t, err
	}
	if containsOnlySpaces(text) {
		err = fmt.Errorf(ERR_NO_TEXT)
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
		return nil, fmt.Errorf("%v : %v", ERR_GETDATA, err)
	}
	listDesc := sortByPriorityDesc(list)
	return listDesc, nil
}

func Add(text Task, priority int, user users.User) (td Todo, err error) {
	td.Task = text
	td.Priority = priority
	err = store.AddTodoDb(td, user)
	if err != nil {
		return td, err
	}
	return td, nil
}

func Delete(id int) (td Todo, err error) {
	td.Id = id
	err = store.DeleteTodoDb(td)
	if err != nil {
		return td, err
	}
	return td, nil
}

func Modify(text Task, id int, priority int) (td Todo, err error) {
	td.Id = id
	td.Task = text
	td.Priority = priority
	err = store.ModifyTodoDb(td)
	if err != nil {
		return td, err
	}
	return td, nil
}
