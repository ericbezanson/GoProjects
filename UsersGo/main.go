package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// User model
type User struct {
	ID    string `gorm:"primaryKey"`
	Name  string
	Email string
}

var db *gorm.DB

func main() {
	// Initialize the database
	var err error
	db, err = gorm.Open(sqlite.Open("./data/test.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	// Migrate the schema
	db.AutoMigrate(&User{})

	// Initialize the router
	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/", handleMainPage).Methods("GET")
	r.HandleFunc("/all-users", handleAllUsers).Methods("GET")
	r.HandleFunc("/users/{uid}", handleUserByID).Methods("GET")
	r.HandleFunc("/users", handleCreateUser).Methods("POST")

	// Start the server
	fmt.Println("Server running on http://localhost:3000")
	http.ListenAndServe(":3000", r)
}

func handleMainPage(w http.ResponseWriter, r *http.Request) {
	// Simple HTML form
	html := `
		<!DOCTYPE html>
		<html>
		<head><title>Create User</title></head>
		<body>
			<h1>Create a New User</h1>
			<form action="/users" method="POST">
				<label for="name">Name:</label><br>
				<input type="text" id="name" name="name"><br><br>
				<label for="email">Email:</label><br>
				<input type="email" id="email" name="email"><br><br>
				<input type="submit" value="Create User">
			</form>
			<a href="/all-users">View All Users</a>
		</body>
		</html>
	`
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func handleAllUsers(w http.ResponseWriter, r *http.Request) {
	var users []User
	db.Find(&users)

	// Build HTML to display all users
	html := "<h1>All Users</h1><table border='1'><tr><th>UID</th><th>Name</th><th>Email</th></tr>"
	for _, user := range users {
		html += fmt.Sprintf("<tr><td><a href='/users/%s'>%s</a></td><td>%s</td><td>%s</td></tr>",
			user.ID, user.ID, user.Name, user.Email)
	}
	html += "</table><br><a href='/'>Back to Create User</a>"

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func handleUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid := vars["uid"]

	var user User
	if result := db.First(&user, "id = ?", uid); result.Error != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Display user details
	html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head><title>User Details</title></head>
		<body>
			<h1>User Details</h1>
			<p><strong>UID:</strong> %s</p>
			<p><strong>Name:</strong> %s</p>
			<p><strong>Email:</strong> %s</p>
			<br><a href="/all-users">Back to All Users</a><br>
			<a href="/">Back to Create User</a>
		</body>
		</html>
	`, user.ID, user.Name, user.Email)

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func handleCreateUser(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	r.ParseForm()
	name := r.FormValue("name")
	email := r.FormValue("email")

	// Create a new user with a UUID as ID
	user := User{
		ID:    generateUUID(),
		Name:  name,
		Email: email,
	}

	// Save the user to the database
	if result := db.Create(&user); result.Error != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Redirect back to the main page or show a success message
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Function to generate a simple UUID (replace with a real UUID generator for production)
// Function to generate a proper UUID
func generateUUID() string {
	return uuid.New().String()
}
