package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

// SocialMediaCheck represents a social media platform to check
type SocialMediaCheck struct {
	URL      string // URL template with %s for username
	Username string // Username to check
}

// Result holds the scanning result for each platform
type Result struct {
	Platform string `json:"platform"`        // Platform name (derived from URL)
	URL      string `json:"url"`             // Full URL checked
	Exists   bool   `json:"exists"`          // Whether the profile exists
	Error    string `json:"error,omitempty"` // Error message if any
}

// loadPlatformsFromFile reads social media URL templates from a text file
func loadPlatformsFromFile(filename string, username string) ([]SocialMediaCheck, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	var platforms []SocialMediaCheck
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") { // Skip empty lines and comments
			parts := strings.SplitN(line, ": ", 2)
			if len(parts) == 2 {
				url := strings.TrimSpace(parts[1])
				platforms = append(platforms, SocialMediaCheck{
					URL:      url,
					Username: username,
				})
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return platforms, nil
}

// checkSocialMedia scans multiple platforms concurrently
func checkSocialMedia(username string) []Result {
	// Load platforms from file
	platforms, err := loadPlatformsFromFile("socials.txt", username)
	if err != nil {
		return []Result{{Platform: "Error", Error: err.Error()}}
	}

	// Channel to collect results from goroutines
	resultChan := make(chan Result, len(platforms))
	var wg sync.WaitGroup

	// HTTP client with timeout
	client := &http.Client{Timeout: 10 * time.Second}

	// Launch a goroutine for each platform
	for _, platform := range platforms {
		wg.Add(1)
		go func(p SocialMediaCheck) {
			defer wg.Done()

			// Construct the full URL with the username
			url := fmt.Sprintf(p.URL, p.Username)
			resp, err := client.Get(url)
			// Extract platform name from URL (simple approach)
			platformName := strings.Split(strings.Split(url, "://")[1], ".")[0]
			result := Result{Platform: platformName, URL: url}

			if err != nil {
				result.Error = err.Error()
			} else {
				defer resp.Body.Close()
				// Simple check: if status is 200, assume profile exists
				result.Exists = resp.StatusCode == http.StatusOK
			}

			// Send result to channel
			resultChan <- result
		}(platform)
	}

	// Close the channel once all goroutines are done
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results from the channel
	results := make([]Result, 0, len(platforms))
	for result := range resultChan {
		results = append(results, result)
	}

	return results
}

// scanHandler handles the API request to scan a username
func scanHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	// Perform the scan
	results := checkSocialMedia(username)

	// Set response headers and encode results as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func main() {
	// Initialize the router
	router := mux.NewRouter()

	// Define the API endpoint
	router.HandleFunc("/scan/{username}", scanHandler).Methods("GET")

	// Start the server
	fmt.Println("Server is running on port 8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}
