package main

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
)

var (
	proxyTarget = os.Getenv("PROXY_TARGET")
	logLevel    = os.Getenv("LOG_LEVEL")
)

func checkRequest(w http.ResponseWriter, r *http.Request) bool {
	// Check if the request is to /rest/v3/short-urls
	if !strings.Contains(r.URL.Path, "/rest/v3/short-urls") {
		if logLevel == "debug" {
			log.Printf("Invalid URL, end point not supported: %s\n", r.URL.Path)
		}
		return false
	}

	return true
}

func checkGetRequest(w http.ResponseWriter, r *http.Request) bool {
	// Check if the URL contains a code
	code := strings.Split(r.URL.Path, "/rest/v3/short-urls/")
	if len(code) < 2 {
		if logLevel == "debug" {
			log.Printf("Invalid URL, missing short link, : %s\n", r.URL.Path)
		}
		return false
	}

	// Check if the code is empty
	shortlink := code[len(code)-1]
	if shortlink == "" {
		if logLevel == "debug" {
			log.Printf("Invalid URL, short link is empty: %s\n", r.URL.Path)
		}
		return false
	}

	return true
}

func checkBody(r *http.Request) (bool, []byte, string) {
	// Check if the content type is application/json
	if r.Header.Get("Content-Type") != "application/json" {
		if logLevel == "debug" {
			log.Printf("Invalid content type: %s\n", r.Header.Get("Content-Type"))
		}
		return false, nil, ""
	}
	if r.Body == nil {
		if logLevel == "debug" {
			log.Printf("Empty body\n")
		}
		return false, nil, ""
	}
	// Read the body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		if logLevel == "debug" {
			log.Printf("Error reading body\n")
		}
		return false, nil, ""
	}

	var parsedBody map[string]interface{}
	json.Unmarshal([]byte(body), &parsedBody)

	if parsedBody == nil {
		if logLevel == "debug" {
			log.Printf("Invalid body, JSON could not be parsed\n")
		}
		return false, nil, ""
	}

	// Check if the body contains the key "longUrl"
	if _, ok := parsedBody["longUrl"]; !ok {
		if logLevel == "debug" {
			log.Printf("Invalid body, missing longUrl key\n")
		}
		return false, nil, ""
	}

	// Check if the longUrl is a string.
	longUrl := parsedBody["longUrl"]
	if longUrl == nil {
		if logLevel == "debug" {
			log.Printf("Invalid body, longUrl is nil\n")
		}
		return false, nil, ""
	}
	if reflect.TypeOf(longUrl).Kind() != reflect.String {
		if logLevel == "debug" {
			log.Printf("Invalid body, longUrl is not a string\n")
		}
		return false, nil, ""
	}

	longUrlS := longUrl.(string)

	return true, body, longUrlS
}

func checkLongUrl(longUrl string) (bool, map[string]interface{}) {
	// Check if the longUrl is valid.
	if !strings.Contains(longUrl, "import/") {
		if logLevel == "debug" {
			log.Printf("Invalid body, longUrl does not contain 'import/'\n")
		}
		return false, nil
	}

	urlParts := strings.Split(longUrl, "import/")
	if len(urlParts) < 2 {
		if logLevel == "debug" {
			log.Printf("Invalid body, longUrl does not contain part after import/\n")
		}
		return false, nil
	}

	base64Str := urlParts[len(urlParts)-1]

	shortcut, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		if logLevel == "debug" {
			log.Printf("Invalid body, base64 decode failed\n")
		}
		return false, nil
	}

	var jsonMap map[string]interface{}
	unmarshErr := json.Unmarshal([]byte(shortcut), &jsonMap)

	if unmarshErr != nil {
		if logLevel == "debug" {
			log.Printf("Invalid body, base64v JSON could not be parsed\n")
		}
		return false, nil
	}

	return true, jsonMap
}

func checkShortcut(jsonMap map[string]interface{}) (bool, string) {
	if jsonMap == nil {
		if logLevel == "debug" {
			log.Printf("Invalid body, JSON could not be parsed\n")
		}
		return false, ""
	}

	// Check the type of the shortcut (e.g. "ShortcutLocation" or "ShortcutRoute")
	shortcutType := jsonMap["type"]

	if shortcutType == nil {
		if logLevel == "debug" {
			log.Printf("Invalid body, missing type key\n")
		}
		return false, ""
	}

	if shortcutType != "ShortcutLocation" && shortcutType != "ShortcutRoute" {
		if logLevel == "debug" {
			log.Printf("Invalid body, invalid type key\n")
		}
		return false, ""
	}

	if _, ok := jsonMap["id"]; !ok {
		if logLevel == "debug" {
			log.Printf("Invalid body, missing id key\n")
		}
		return false, ""
	}

	if _, ok := jsonMap["name"]; !ok {
		if logLevel == "debug" {
			log.Printf("Invalid body, missing name key\n")
		}
		return false, ""
	}

	return true, shortcutType.(string)
}

func checkLocationShortcut(shortcut map[string]interface{}) bool {
	if len(shortcut) > 4 {
		if logLevel == "debug" {
			log.Printf("Invalid body, too many keys\n")
		}
		return false
	}

	if _, ok := shortcut["waypoint"]; !ok {
		if logLevel == "debug" {
			log.Printf("Invalid body, missing waypoint key\n")
		}
		return false
	}

	return true
}

func checkRouteShortcut(shortcut map[string]interface{}) bool {
	if len(shortcut) > 6 {
		if logLevel == "debug" {
			log.Printf("Invalid body, too many keys\n")
		}
		return false
	}

	if _, ok := shortcut["waypoints"]; !ok {
		if logLevel == "debug" {
			log.Printf("Invalid body, missing waypoints key\n")
		}
		return false
	}

	if _, ok := shortcut["routeTimeText"]; !ok {
		if logLevel == "debug" {
			log.Printf("Invalid body, missing routeTimeText key\n")
		}
		return false
	}

	if _, ok := shortcut["routeLengthText"]; !ok {
		if logLevel == "debug" {
			log.Printf("Invalid body, missing routeLengthText key\n")
		}
		return false
	}

	return true
}

func checkAndProxy(w http.ResponseWriter, r *http.Request) {
	if logLevel == "debug" {
		log.Printf("Request: %s %s\n", r.Method, r.URL.Path)
	}

	// Health check endpoint
	if r.URL.Path == "/health" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if !checkRequest(w, r) {
		http.Error(w, "Invalid", http.StatusBadRequest)
		return
	}

	// If this is a GET request, we only want to allow /rest/v3/short-urls/{code}
	if r.Method == http.MethodGet {
		if !checkGetRequest(w, r) {
			http.Error(w, "Invalid", http.StatusBadRequest)
			return
		}
		proxy(w, r, nil)
		return
	}

	// If this is a POST request, we need to check the content
	if r.Method == http.MethodPost {
		ok, body, longUrl := checkBody(r)
		if !ok {
			http.Error(w, "Invalid", http.StatusBadRequest)
			return
		}

		ok, jsonMap := checkLongUrl(longUrl)
		if !ok {
			http.Error(w, "Invalid", http.StatusBadRequest)
			return
		}

		ok, shortcutType := checkShortcut(jsonMap)
		if !ok {
			http.Error(w, "Invalid", http.StatusBadRequest)
			return
		}

		if shortcutType == "ShortcutLocation" {
			if !checkLocationShortcut(jsonMap) {
				http.Error(w, "Invalid", http.StatusBadRequest)
				return
			}
		}

		if shortcutType == "ShortcutRoute" {
			if !checkRouteShortcut(jsonMap) {
				http.Error(w, "Invalid", http.StatusBadRequest)
				return
			}
		}

		// If the JSON contains the required keys, proxy the request
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
