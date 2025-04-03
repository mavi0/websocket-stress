package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mavi0/websocket-stress/internal/server"
)

func main() {
	// Get current working directory
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting current directory: ", err)
	}

	// Navigate up two levels to find project root
	// If running from cmd/websocket-stress, this gives us the project root
	rootDir := workingDir
	if filepath.Base(workingDir) == "websocket-stress" && filepath.Base(filepath.Dir(workingDir)) == "cmd" {
		rootDir = filepath.Join(workingDir, "../..")
	}

	// Convert to absolute path and clean it
	rootDir, err = filepath.Abs(rootDir)
	if err != nil {
		log.Fatal("Error getting absolute path: ", err)
	}

	// Initialize the websocket manager
	manager := server.NewManager()

	// Start the manager's goroutines
	go manager.Start()
	go manager.SendUpdates()

	// Verify template path exists
	templatePath := filepath.Join(rootDir, "web/templates/index.html")
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		log.Printf("WARNING: Template file not found at: %s", templatePath)
	} else {
		log.Printf("Template file found at: %s", templatePath)
	}

	// Set up HTTP routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		server.ServeHome(w, r, rootDir)
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		server.ServeWs(manager, w, r)
	})

	// Serve static files from the correct location
	staticDir := http.FileServer(http.Dir(filepath.Join(rootDir, "web/static")))
	http.Handle("/static/", http.StripPrefix("/static/", staticDir))

	// Start the server on all interfaces
	port := "0.0.0.0:8088"
	log.Printf("Server starting on %s", port)
	log.Printf("Using root directory: %s", rootDir)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
