package cmd

import (
	"fmt"
)

func Help() {
	fmt.Println("Mytodolist en CLI !")
	fmt.Println("Usage: mytodolist <command> [arguments]")
	fmt.Println("")
	fmt.Println("Commandes disponibles:")
	fmt.Println("  get          		Affiche vos tâches")
	fmt.Println("  add          		Ajouter une nouvelle tâche")
	fmt.Println("  delete		Supprime une tâche existante")
	fmt.Println("  modify		Modifie une tâche existante")
	fmt.Println("  logout                Déconnecte l'utilisateur")
	fmt.Println("  help                	Affiche ce message")
	fmt.Println("")
}
