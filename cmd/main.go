package main

import (
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

	log.Printf("Received request for %d: %s", statusCode, statusText)

	// Write the header
	w.WriteHeader(statusCode)

	// Write the body
	fmt.Fprintf(w, "%d %s\n", statusCode, statusText)

	fmt.Fprintln(w, "\nHeaders:")
	log.Println("Headers:")
	for name, values := range r.Header {
		for _, value := range values {
			fmt.Fprintf(w, "%s: %s\n", name, value)
			log.Printf("%s: %s", name, value)
		}
	}
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
