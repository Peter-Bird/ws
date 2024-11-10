// File: handler.go
package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// Handler function to use reflection to choose the appropriate service
func Handler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	log.Printf("Received request")

	// Extract the last part of the URL path as the service name
	pathParts := strings.Split(r.URL.Path, "/")
	serviceType := pathParts[len(pathParts)-1] // Get the endpoint name
	log.Printf("Service type requested: %s", serviceType)

	var input map[string]interface{}

	// Support both GET and POST methods
	if r.Method == http.MethodGet {
		input = make(map[string]interface{})
		for key, values := range r.URL.Query() {
			if len(values) > 0 {
				input[key] = values[0]
			}
		}
		log.Printf("Request payload from query parameters: %v", input)
	} else if r.Method == http.MethodPost {
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			log.Printf("Error decoding JSON: %v", err)
			return
		}
		log.Printf("Request payload from body: %v", input)
	} else {
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		log.Printf("Unsupported method: %s", r.Method)
		return
	}

	service, exists := GetService(serviceType)
	if !exists {
		http.Error(w, "Unknown service type", http.StatusBadRequest)
		log.Printf("Unknown service type: %s", serviceType)
		return
	}

	result, err := service.Process(input)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error processing request: %v", err), http.StatusInternalServerError)
		log.Printf("Error processing request with service %s: %v", serviceType, err)
		return
	}
	log.Printf("Service %s processed request successfully", serviceType)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("Error encoding response: %v", err)
	}

	endTime := time.Now()
	log.Printf("Request completed, duration: %v", endTime.Sub(startTime))
}
