package cmd

import (
	"fmt"
)

func Help(cmd string) {
	switch cmd {
	case "help":
		fmt.Println("My TodoList en CLI !")
		fmt.Println("Usage: mytodolist <command> [arguments]")
		fmt.Println("\nCommandes disponibles:")
		fmt.Println("  get          		Affiche vos tâches")
		fmt.Println("  add          		Ajouter une nouvelle tâche")
		fmt.Println("  delete		Supprime une tâche existante")
		fmt.Println("  modify		Modifie une tâche existante")
		fmt.Println("  signin                S'inscrire")
		fmt.Println("  login                 Se connecter")
		fmt.Println("  logout                Se déconnecter")
		fmt.Println("  help                	Affiche ce message")
	case "modify":
		fmt.Println(`Usage: mytodolist modify <id> "text" <priority>`)
		fmt.Println("\nid :             l'identifiant numérique du todo que vous souhaitez modifier")
		fmt.Println("text :           une chaîne de caractère décrivant votre todo")
		fmt.Println("priority :       le niveau de priorité de votre tâche entre 1 et 3, du moins au plus urgent")
	case "add":
		fmt.Println(`Usage: mytodolist add "text" <priority>`)
		fmt.Println("\ntext :           une chaîne de caractère décrivant votre nouveau todo")
		fmt.Println("priority :       le niveau de priorité de votre tâche entre 1 et 3, du moins au plus urgent")
	case "delete":
		fmt.Println("Usage: mytodolist delete <id>")
		fmt.Println("\nid :             l'identifiant numérique du todo que vous souhaitez supprimer")
	}
}
