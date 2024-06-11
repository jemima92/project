package handlers

import (
	"myproject/db"
	"net/http"
	"text/template"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var tpl *template.Template

func InitTemplates() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

type Post struct {
	ID       int
	Username string
	Title    string
	Body     string
	Date     string
}

func TestDBConnection(w http.ResponseWriter, r *http.Request) {
	err := db.DB.Ping()
	if err != nil {
		http.Error(w, "Failed to connect to the database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Successfully connected to the database"))
}

func Index(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT posts.id, users.username, posts.title, posts.body, posts.date FROM posts JOIN users ON posts.user_id = users.id ORDER BY posts.date DESC")
	if err != nil {
		http.Error(w, "Internal server error: unable to query database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Username, &post.Title, &post.Body, &post.Date)
		if err != nil {
			http.Error(w, "Internal server error: unable to scan row", http.StatusInternalServerError)
			return
		}
		post.Date = formatDate(post.Date)
		posts = append(posts, post)
	}

	if err = tpl.ExecuteTemplate(w, "index.gohtml", posts); err != nil {
		http.Error(w, "Internal server error: unable to execute template", http.StatusInternalServerError)
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Internal server error: unable to hash password", http.StatusInternalServerError)
			return
		}

		_, err = db.DB.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, hashedPassword)
		if err != nil {
			http.Error(w, "Internal server error: unable to insert user", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if err := tpl.ExecuteTemplate(w, "register.gohtml", nil); err != nil {
		http.Error(w, "Internal server error: unable to execute template", http.StatusInternalServerError)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var hashedPassword string
		err := db.DB.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&hashedPassword)
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "username",
			Value: username,
			Path:  "/",
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if err := tpl.ExecuteTemplate(w, "login.gohtml", nil); err != nil {
		http.Error(w, "Internal server error: unable to execute template", http.StatusInternalServerError)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "username",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func NewPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		body := r.FormValue("body")

		cookie, err := r.Cookie("username")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		username := cookie.Value
		var userID int
		err = db.DB.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&userID)
		if err != nil {
			http.Error(w, "User not found", http.StatusInternalServerError)
			return
		}

		_, err = db.DB.Exec("INSERT INTO posts (user_id, title, body) VALUES (?, ?, ?)", userID, title, body)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tpl.ExecuteTemplate(w, "newpost.gohtml", nil)
}
func Categories(w http.ResponseWriter, r *http.Request) {
	if err := tpl.ExecuteTemplate(w, "categories.gohtml", nil); err != nil {
		http.Error(w, "Internal server error: unable to execute template", http.StatusInternalServerError)
	}
}
func formatDate(date string) string {
	t, err := time.Parse("2006-01-02 15:04:05", date)
	if err != nil {
		return date
	}
	return t.Format("Jan 2, 2006 at 3:04pm")
}
