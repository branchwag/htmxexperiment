package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

const divHTML = `<div id="response">Button was clicked!</div>`
const formSubmitHTML = `<div id="response">Form submitted!</div>`

func clickedHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost{
		t, err := template.New("response").Parse(divHTML)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		if err := t.Execute(w, nil); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, filepath.Join(".", "index.html"))
}

func otherpageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/otherpage" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, filepath.Join(".", "otherpage.html"))
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost{
		t, err := template.New("response").Parse(formSubmitHTML)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		if err := t.Execute(w, nil); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/clicked", clickedHandler)
	http.HandleFunc("/otherpage", otherpageHandler)
	http.HandleFunc("/submit", submitHandler)

	fmt.Println("Starting server on :4242")
	if err := http.ListenAndServe(":4242", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}