package main

import "fmt"

type Todo struct {
	done bool
	text string
}

type MyTodo struct {
	todos []Todo
}

//créer, modifier, supprimer todo
func (mt *MyTodo) add(t Todo) error {
	if t.text == "" {
		return fmt.Errorf("Pas de texte renseigné !")
	}
	mt.todos = append(mt.todos, t)
	return nil
}

func (mt *MyTodo) delete(text string) error {
	for i, t := range mt.todos {
		if t.text == text {
			//Supprime l'élément visé sans changer l'ordre
			mt.todos = append(mt.todos[:i], mt.todos[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Todo '%s' introuvable", text)
}

func (mt *MyTodo) modif(oldText, newText string) error {
	for i, t := range mt.todos {
		if t.text == oldText {
			mt.todos[i].text = newText
			return nil
		}
	}
	return fmt.Errorf("Todo '%s' introuvable", oldText)
}

func main() {
	a := Todo{false, "Ma 1ère tâche"}
	b := Todo{false, "Ma 2ème tâche"}
	todos := []Todo{}
	myList := MyTodo{todos}

	myList.add(a)
	myList.add(b)
	fmt.Println(myList)

	err := myList.modif("Ma 1ère tâche", "1ère tâche modifiée")
	if err != nil {
		fmt.Println("Error :", err)
	} else {
		fmt.Println("Nouvelle liste :", myList)
	}

	newTodo := Todo{false, ""}
	errAdd := myList.add(newTodo)
	if errAdd != nil {
		fmt.Println("Error :", errAdd)
	} else {
		fmt.Println("Todo ajouté !", myList)
	}

	errDel := myList.delete("Ma 2ème tâche")
	if errDel != nil {
		fmt.Println("Error :", errDel)
	} else {
		fmt.Println("Todo supprimé !", myList)
	}
}
