package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
)

const divHTML = `<div id="response">Button was clicked!</div>`
const formSubmitHTML = `<div id="response">Form submitted!</div>`
const fooHTML = `<div id="response">You fooed it up!</div>`

//loads the environment vars
func LoadEnv(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		os.Setenv(key, value)
	}
	return scanner.Err()
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

func testHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		t, err := template.New("testresponse").Parse(fooHTML)
		if err != nil {
			http.Error(w, "Internal Foo Error", http.StatusInternalServerError)
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

		err := LoadEnv(".env")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Error loading .env file:", err)
			return
		}

		server := os.Getenv("DB_SERVER")
		//port := os.Getenv("DB_PORT")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		database := os.Getenv("DB_NAME")

		connString := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s;", server, user, password, database)

		db, err := sql.Open("sqlserver", connString)
		if err != nil {
			http.Error(w, "Internal Server Error DB", http.StatusInternalServerError)
			log.Println("Database connection error:", err)
			return
		}
		defer db.Close()

		// Create table if it doesn't exist
		sqlStmt := `
		IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='users' AND xtype='U')
		CREATE TABLE users (
			id INT IDENTITY(1,1) PRIMARY KEY,
			name NVARCHAR(100),
			email NVARCHAR(100)
		);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Error creating table:", err)
			return
		}

		// Insert the form data into the database using named parameters
		_, err = db.Exec("INSERT INTO users (name, email) VALUES (@name, @email)",
			sql.Named("name", name),
			sql.Named("email", email))
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Error inserting data:", err)
			return
		}

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


func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/clicked", clickedHandler)
	http.HandleFunc("/otherpage", otherpageHandler)
	http.HandleFunc("/submit", submitHandler)
	http.HandleFunc("/test", testHandler)

	fmt.Println("Starting server on :4242")
	if err := http.ListenAndServe(":4242", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
