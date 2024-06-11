package main

import (
	"fmt"
	"log"
	"myproject/db"
	"myproject/handlers"
	"net/http"
	"os/exec"
	"runtime"
)

func main() {
	// Initialize the database
	db.InitDB()

	// Initialize templates
	handlers.InitTemplates()

	// Set up routes
	http.HandleFunc("/", handlers.Index)
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/testdb", handlers.TestDBConnection) // Add this line
	http.HandleFunc("/categories", handlers.Categories)   // Add this line

	// Open the browser after starting the server
	go func() {
		openBrowser("http://localhost:8080")
	}()

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// openBrowser tries to open the URL in a web browser
func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin": // macOS
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		log.Println("Error opening browser:", err)
	}
}
