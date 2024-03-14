package main

import "fmt"

type Todo struct {
	Done bool
	Text string
}

type MyTodoList struct {
	TodoList []Todo
}

var todos = []Todo{}
var myList = MyTodoList{todos}

//créer, modifier, supprimer todo
func (mt *MyTodoList) add(todo Todo) error {
	existingTodo := false
	if todo.Text == "" {
		return fmt.Errorf("Pas de texte renseigné !")
	} else {
		for _, t := range myList.TodoList {
			if t.Text == todo.Text {
				existingTodo = true
				return fmt.Errorf("Todo existant !")
			}
		}
		if existingTodo == false {
			mt.TodoList = append(mt.TodoList, todo)
		}
	}
	return nil
}

func (mt *MyTodoList) delete(todo Todo) error {
	for i, t := range mt.TodoList {
		if t.Text == todo.Text {
			//Supprime l'élément visé sans changer l'ordre
			mt.TodoList = append(mt.TodoList[:i], mt.TodoList[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Todo '%s' introuvable", todo.Text)
}

func (mt *MyTodoList) modif(oldText, newText string) error {
	for i, t := range mt.TodoList {
		if t.Text == oldText {
			mt.TodoList[i].Text = newText
			return nil
		}
	}
	return fmt.Errorf("Todo '%s' introuvable", oldText)
}

func main() {
	a := Todo{false, "Ma 1ère tâche"}
	b := Todo{true, "Ma 2ème tâche"}
	todos := []Todo{}
	myList := MyTodoList{todos}

	myList.add(a)
	myList.add(b)
	fmt.Println(myList)

	err := myList.modif("Ma 1ère tâche", "1ère tâche modifiée")
	if err != nil {
		fmt.Println("Error :", err)
	} else {
		fmt.Println("Nouvelle liste :", myList)
	}

	newTodo := Todo{false, "C moi"}
	errAdd := myList.add(newTodo)
	if errAdd != nil {
		fmt.Println("Error :", errAdd)
	} else {
		fmt.Println("Todo ajouté !", myList)
	}

	errDel := myList.delete(b)
	if errDel != nil {
		fmt.Println("Error :", errDel)
	} else {
		fmt.Println("Todo supprimé !", myList)
	}
}
