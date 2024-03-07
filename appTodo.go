package main

import (
	"net/http"
)

func main() {
	// handle route using handler function
	fs := http.FileServer(http.Dir("./ihm"))
	http.Handle("/", fs)

	// listen to port
	http.ListenAndServe(":5050", nil)
}
