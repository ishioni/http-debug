package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func statusHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the path, removing the leading slash
	path := strings.TrimPrefix(r.URL.Path, "/")

	// If path is empty, return 200 OK with a help message
	if path == "" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Welcome to http-debug! usage: /<status-code> (e.g., /200, /404, /500)")
		return
	}

	// Try to parse the path as an integer
	statusCode, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid status code: %s", path), http.StatusBadRequest)
		return
	}

	statusText := http.StatusText(statusCode)
	if statusText == "" {
		statusText = "Unknown Status Code"
	}

	data := struct {
		StatusCode int         `json:"statusCode"`
		StatusText string      `json:"statusText"`
		Headers    http.Header `json:"headers"`
	}{
		StatusCode: statusCode,
		StatusText: statusText,
		Headers:    r.Header,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("ERROR: could not marshal json: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Printf("%s", jsonData)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonData)
}

func readyzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK")
}

func main() {
	http.HandleFunc("/", statusHandler)
	http.HandleFunc("/readyz", readyzHandler)

	port := ":8080"
	log.Printf("Starting http-debug server on port %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
