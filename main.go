package main

// Author:Mateusz Kędziora https://mkedziora.pl

import (
	"encoding/json"
	"fmt"
	"log"
	"mkedziora/fast-links/kvstore"
	"net/http"
	"os"
	"strconv"
	"time"
)

var store *kvstore.KVStore

func init() {
	var err error
	filePath := os.Getenv("DATA_FILE_PATH")
	store, err = kvstore.NewKVStore(filePath)
	if err != nil {
		log.Fatalf("Failed to create KVStore: %v", err)
	}
}

func setHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		Data   string `json:"value"`
		ID     string `json:"id,omitempty"`
		Expire string `json:"expire,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// mu.Lock()
	if input.ID == "" {
		input.ID = "default"
	}

	expireString, err := strconv.Atoi(input.Expire)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	expireSeconds := 30
	if expireString > 0 {
		if expireString > 86400 {
			http.Error(w, fmt.Sprintf("Maksymalny czas wygaśnięcia to %d sekund", 86400), http.StatusBadRequest)
			return
		}
		expireSeconds = expireString
	}

	expirationTime := time.Now().Add(time.Duration(expireSeconds) * time.Second)

	store.SetWithTimestamp("key_"+input.ID, input.Data, expirationTime)
	// dataMap[input.ID] = DataEntry{
	// 	Value:     input.Data,
	// 	Timestamp: expirationTime,
	// }
	// mu.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		id = "default"
	}

	// mu.Lock()
	// entry, exists := dataMap[id]
	// mu.Unlock()

	entry, exists := store.Get("key_" + id)

	if !exists || time.Since(entry.Timestamp) > 30*time.Second {
		store.Delete("key_" + id)
		http.Error(w, "No data available", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"data": entry.Value})
}

func getUrlHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		id = "default"
	}

	entry, exists := store.Get("key_" + id)
	if !exists || time.Since(entry.Timestamp) > 30*time.Second {
		store.Delete("key_" + id)
		http.Error(w, "No data available", http.StatusNotFound)
		return
	}

	w.Header().Set("Referrer-Policy", "no-referrer")
	http.Redirect(w, r, entry.Value, http.StatusFound)
}

func main() {

	// Obsługa endpointów API
	http.HandleFunc("/api/set", setHandler)
	http.HandleFunc("/api/get", getHandler)
	http.HandleFunc("/api/url", getUrlHandler)

	// Serwowanie plików statycznych
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "public/index.html")
		} else {
			http.NotFound(w, r)
		}
	})

	http.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/app.html")
	})

	http.HandleFunc("/css/index.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/css/index.css")
	})

	http.HandleFunc("/css/app.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/css/app.css")
	})

	http.HandleFunc("/js/index.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/js/index.js")
	})

	http.HandleFunc("/js/app.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/js/app.js")
	})

	go func() {
		for {

			for id, timestamp := range store.GetAll() {
				if time.Now().After(timestamp.Timestamp) {
					store.Delete(id)
					break
				}
			}
			time.Sleep(5 * time.Second)
		}
	}()

	// Uruchomienie serwera
	http.ListenAndServe(":8080", nil)
}
