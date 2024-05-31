package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	proxyTarget = os.Getenv("PROXY_TARGET")
	logLevel    = os.Getenv("LOG_LEVEL")
)

func checkAndProxy(w http.ResponseWriter, r *http.Request) {
	if logLevel == "debug" {
		log.Printf("Request: %s %s\n", r.Method, r.URL.Path)
	}

	// Check if the request is to /rest/v3/short-urls
	if !strings.Contains(r.URL.Path, "/rest/v3/short-urls") {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	// If this is a GET request, we don't need to check anything
	if r.Method == http.MethodGet {
		proxy(w, r, nil)
		return
	}

	// If this is a POST request, we need to check the content type
	if r.Method == http.MethodPost {
		// Check if the content type is application/json
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Invalid content type", http.StatusBadRequest)
			return
		}
		if r.Body == nil {
			http.Error(w, "Empty body", http.StatusBadRequest)
			return
		}
		// Read the body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading body", http.StatusInternalServerError)
			return
		}
		// Check if the body contains the key "longUrl"
		if !strings.Contains(string(body), "longUrl") {
			http.Error(w, "Invalid body", http.StatusBadRequest)
			return
		}
		
		// TODO: Make more checks on the body, e.g. validate the JSON structure
		// and whether the contained code can be parsed as a valid shortcut.

		proxy(w, r, body)
		return
	}
}

func proxy(w http.ResponseWriter, r *http.Request, body []byte) {
	if logLevel == "debug" {
		log.Printf("Proxying request: %s %s\n", r.Method, r.URL.Path)
	}

	// Create a new request to the target server
	targetURL := proxyTarget + r.URL.Path
	proxyReq, err := http.NewRequest(r.Method, targetURL, strings.NewReader(string(body)))
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}
	if logLevel == "debug" {
		log.Printf("Proxied URL: %s\n", targetURL)
		log.Printf("Proxied method: %s\n", r.Method)
	}

	// Copy headers from the original request to the proxy request
	for key, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(key, value)
			if logLevel == "debug" {
				log.Printf("Header: %s: %s\n", key, value)
			}
		}
	}

	// Perform the request to the target server
	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, "Error performing request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy the response from the target server to the original response
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)

	if logLevel == "debug" {
		log.Printf("Proxied response: %d\n", resp.StatusCode)
	}
}

func main() {
	// Create a new HTTP server with the handleRequest function as the handler
	server := http.Server{
		Addr:    ":8000",
		Handler: http.HandlerFunc(checkAndProxy),
	}

	// Start the server and log any errors
	log.Println("Starting proxy server on :8000")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting proxy server: ", err)
	}
}
