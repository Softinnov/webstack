package main

import (
	"net/http"
)

func main() {
	// handle route using handler function
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "ihm/src/app.svelte")
	})

	// listen to port
	http.ListenAndServe(":5050", nil)
}
