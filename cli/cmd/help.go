package cmd

import (
	"fmt"
)

func Help() {
	fmt.Println("Welcome to mytodolist CLI app!")
	fmt.Println("Usage: mytodolist <command> [arguments]")
	fmt.Println("")
	fmt.Println("Available commands:")
	fmt.Println("  get          		List all tasks")
	fmt.Println("  help                	Show this help message")
	fmt.Println("")
}
