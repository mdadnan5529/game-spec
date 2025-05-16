// Simple Go backend for managing gaming hardware and games.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

// Data structures for storing gaming hardware and games.
type Hardware struct {
	Type string `json:"type"` // e.g., Laptop, PC, Console
	Name string `json:"name"`
}

type Game struct {
	Title string `json:"title"`
}

var (
	httpAddr = flag.String("addr", ":8080", "HTTP listen address")
)

func main() {
	flag.Parse()

	// Simple in-memory storage
	hardwareList := []Hardware{}
	gameList := []Game{}
	var mu sync.Mutex // Mutex for protecting access to hardwareList and gameList

	// Handler for adding new hardware
	http.HandleFunc("/add-hardware", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}

		var newHardware Hardware
		err = json.Unmarshal(body, &newHardware)
		if err != nil {
			http.Error(w, "Error parsing JSON", http.StatusBadRequest)
			return
		}

		mu.Lock()
		hardwareList = append(hardwareList, newHardware)
		mu.Unlock()

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "Hardware added successfully")
	})

	// Handler for adding new games
	http.HandleFunc("/add-game", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}

		var newGame Game
		err = json.Unmarshal(body, &newGame)
		if err != nil {
			http.Error(w, "Error parsing JSON", http.StatusBadRequest)
			return
		}

		mu.Lock()
		gameList = append(gameList, newGame)
		mu.Unlock()

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "Game added successfully")
	})

	// Handler for viewing hardware
	http.HandleFunc("/view-hardware", func(w http.ResponseWriter, r *http.Request) {
		mu.RLock()
		defer mu.RUnlock()

		json.NewEncoder(w).Encode(hardwareList)
	})

	// Handler for viewing games
	http.HandleFunc("/view-games", func(w http.ResponseWriter, r *http.Request) {
		mu.RLock()
		defer mu.RUnlock()

		json.NewEncoder(w).Encode(gameList)
	})

	// Handler for serving the index.html file
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "index.html")
	})

	log.Printf("serving http://%s", *httpAddr)
	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}
