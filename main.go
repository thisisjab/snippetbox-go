package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Hello from Snippetbox"))
}

func showSnippets(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.Header().Set("Allow", "GET")
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.Atoi((r.URL.Query().Get("id")))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Visiting snippet id %d.", id)
}

func createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed."))
		return
	}

	w.Write([]byte("Trying to create a snippet."))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippets/create", createSnippet)
	mux.HandleFunc("/snippets", showSnippets)

	log.Println("Starting server on :4000")

	err := http.ListenAndServe(":4000", mux)

	log.Fatal(err)
}
