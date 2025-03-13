package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

// SocialMediaCheck represents a social media platform to check
type SocialMediaCheck struct {
	Name     string // Name of the platform
	URL      string // URL template with %s for username
	Username string // Username to check
}

// Result holds the scanning result for each platform
type Result struct {
	Platform string `json:"platform"`        // Platform name
	URL      string `json:"url"`             // Full URL checked
	Exists   bool   `json:"exists"`          // Whether the profile exists
	Error    string `json:"error,omitempty"` // Error message if any
}

// checkSocialMedia scans multiple platforms concurrently
func checkSocialMedia(username string) []Result {
	// List of social media platforms to check
	platforms := []SocialMediaCheck{
		{Name: "Twitter", URL: "https://twitter.com/%s", Username: username},
		{Name: "Instagram", URL: "https://instagram.com/%s", Username: username},
		{Name: "Facebook", URL: "https://facebook.com/%s", Username: username},
		{Name: "GitHub", URL: "https://github.com/%s", Username: username},
		{Name: "LinkedIn", URL: "https://linkedin.com/in/%s", Username: username},
		{Name: "Reddit", URL: "https://reddit.com/user/%s", Username: username},
		{Name: "TikTok", URL: "https://tiktok.com/@%s", Username: username},
		{Name: "YouTube", URL: "https://youtube.com/@%s", Username: username},
		{Name: "Pinterest", URL: "https://pinterest.com/%s", Username: username},
		{Name: "Medium", URL: "https://medium.com/@%s", Username: username},
		{Name: "Twitch", URL: "https://twitch.tv/%s", Username: username},
		{Name: "Snapchat", URL: "https://snapchat.com/add/%s", Username: username},
		{Name: "Telegram", URL: "https://t.me/%s", Username: username},
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
			result := Result{Platform: p.Name, URL: url}

			if err != nil {
				result.Error = err.Error()
			} else {
				defer resp.Body.Close()
				// Simple check: if status is 200, assume profile exists
				// Note: Some platforms may require more sophisticated checks
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
