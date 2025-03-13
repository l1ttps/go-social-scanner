# Social Media Username Scanner

## Introduction
This is a simple web application written in Go that uses `mux` to create an API for checking the existence of a username on various social media platforms.

## Features
- Checks for username availability on popular social media platforms such as Twitter, Instagram, Facebook, GitHub, LinkedIn, Reddit, TikTok, YouTube, Pinterest, Medium, Twitch, Snapchat, and Telegram.
- Utilizes Goroutines to perform concurrent checks for faster execution.
- Returns results in JSON format.

## Installation
### System Requirements
- Go 1.16 or later

### How to Run
1. Clone the repository:
   ```sh
   git clone https://github.com/l1ttps/go-social-scanner
   cd go-social-scanner
   ```
2. Install dependencies:
   ```sh
   go mod tidy
   ```
3. Run the server:
   ```sh
   go run main.go
   ```

## API Endpoint
### Username Check
- **Endpoint:** `/scan/{username}`
- **Method:** `GET`
- **Example:**
  ```sh
  curl http://localhost:8080/scan/exampleuser
  ```
- **Sample Response:**
  ```json
  [
    {
      "platform": "Twitter",
      "url": "https://twitter.com/exampleuser",
      "exists": true,
      "error": ""
    },
    {
      "platform": "GitHub",
      "url": "https://github.com/exampleuser",
      "exists": false,
      "error": ""
    }
  ]
  ```

## Project Structure
```
/
├── main.go          # Main source code
├── go.mod           # Dependency management
├── go.sum           # Hash of dependencies
├── README.md        # Documentation
```

## Technologies Used
- Golang
- Gorilla Mux
- Goroutines and Channels

## Notes
- Some platforms may block requests if too many are made in a short period.
- Some social media sites may require authentication or an API token for more accurate checks.
- The results are approximate as some platforms may change their response methods.

## License
This project is released under the MIT License.

## Repository
- GitHub: [go-social-scanner](https://github.com/l1ttps/go-social-scanner)

