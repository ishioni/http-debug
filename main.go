package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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

		// Check if it's a valid status code range (roughly)
		// Go's http package doesn't strictly enforce valid codes in WriteHeader, but it's good to be reasonable.
		// However, for debugging, maybe we want to allow weird ones too?
		// The user asked for "corresponding http code". I'll trust the integer.

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
	})

	port := ":8080"
	log.Printf("Starting http-debug server on port %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
