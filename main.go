package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

const divHTML = `<div id="response">Button was clicked!</div>`
const formSubmitHTML = `<div id="response">Form submitted!</div>`

// In-memory "database"
var (
	userDB = make(map[int]User)
	idSeq  = 1
	mu     sync.Mutex
)

type User struct {
	ID    int
	Name  string
	Email string
}


func clickedHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
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
	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		email := r.FormValue("email")

		// Simulate database insert operation
		mu.Lock()
		defer mu.Unlock()

		userDB[idSeq] = User{
			ID:    idSeq,
			Name:  name,
			Email: email,
		}
		idSeq++

		// Respond with the confirmation message
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

func dbHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/db" {
		http.NotFound(w, r)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	const dbHTML = `
	<html>
	<head><title>In-Memory DB</title></head>
	<body>
		<h1>In-Memory Database Contents</h1>
		<table border="1">
			<tr><th>ID</th><th>Name</th><th>Email</th></tr>
			{{range $id, $user := .}}
				<tr>
					<td>{{ $user.ID }}</td>
					<td>{{ $user.Name }}</td>
					<td>{{ $user.Email }}</td>
				</tr>
			{{end}}
		</table>
	</body>
	</html>`

	t, err := template.New("db").Parse(dbHTML)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if err := t.Execute(w, userDB); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/clicked", clickedHandler)
	http.HandleFunc("/otherpage", otherpageHandler)
	http.HandleFunc("/submit", submitHandler)
	http.HandleFunc("/db", dbHandler)

	fmt.Println("Starting server on :4242")
	if err := http.ListenAndServe(":4242", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
