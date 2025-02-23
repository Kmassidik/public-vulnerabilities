package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// SQL Injection Vulnerability
func vulnerableSQLHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./example.db")
	if (err != nil) {
		log.Fatal(err)
	}
	defer db.Close()

	username := r.URL.Query().Get("username")
	query := "SELECT * FROM users WHERE username = '" + username + "'"
	fmt.Println("Executing Query:", query)

	rows, err := db.Query(query) // Vulnerable to SQL injection
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	fmt.Fprintf(w, "Query executed for user: %s", username)
}

// XSS Vulnerability
func vulnerableXSSHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	fmt.Fprintf(w, "<h1>Welcome, %s!</h1>", name) // Vulnerable to XSS injection
}

// Path Traversal Vulnerability
func vulnerableFileHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("file")
	http.ServeFile(w, r, filePath) // Vulnerable to path traversal
}

// Insecure Deserialization (JSON Injection)
type User struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

func vulnerableJSONHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user) // No validation on input
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Welcome %s! Your role is %s.", user.Name, user.Role)
}

// Insecure Randomness (Predictable Token Generation)
func insecureTokenHandler(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().UnixNano())
	token := fmt.Sprintf("%d", rand.Int()) // Predictable token
	fmt.Fprintf(w, "Generated token: %s", token)
}

// Information Disclosure (Verbose Error Messages)
func verboseErrorHandler(w http.ResponseWriter, r *http.Request) {
	_, err := http.Get("http://invalid-url") // Invalid URL
	if err != nil {
		http.Error(w, fmt.Sprintf("Detailed Error: %v", err), http.StatusInternalServerError) // Reveals too much info
		return
	}
	fmt.Fprintf(w, "Request succeeded")
}

// Improper Error Handling (Silent Failures)
func silentFailureHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		// Silent failure: No logging or feedback
		return
	}
	fmt.Fprintf(w, "Form parsed successfully")
}

func main() {
	http.HandleFunc("/sql", vulnerableSQLHandler)
	http.HandleFunc("/xss", vulnerableXSSHandler)
	http.HandleFunc("/file", vulnerableFileHandler)
	http.HandleFunc("/json", vulnerableJSONHandler)
	http.HandleFunc("/token", insecureTokenHandler)
	http.HandleFunc("/error", verboseErrorHandler)
	http.HandleFunc("/silent", silentFailureHandler)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
