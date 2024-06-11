package db

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./myproject.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	// Create tables if they do not exist
	createTables()
	statement, _ := DB.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT)")
	statement.Exec()
	//statement, _ = DB.Prepare("INSERT INTO people (firstname, lastname) VALUES (?, ?)") => petit test
	//statement.Exec("Nic", "Raboy")
	rows, _ := DB.Query("SELECT id, firstname, lastname FROM people")
	var id int
	var firstname string
	var lastname string
	for rows.Next() {
		rows.Scan(&id, &firstname, &lastname)
		fmt.Println(strconv.Itoa(id) + ": " + firstname + " " + lastname)
	}
}

func createTables() {
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		userid INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);`

	_, err := DB.Exec(createUsersTable)
	if err != nil {
		log.Fatal("Failed to create users table:", err)
	}

	createPostsTable := `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		title TEXT NOT NULL,
		body TEXT NOT NULL,
		date DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(userid)
	);`

	_, err = DB.Exec(createPostsTable)
	if err != nil {
		log.Fatal("Failed to create posts table:", err)
	}
}
