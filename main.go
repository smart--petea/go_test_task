package main

import (
    "github.com/joho/godotenv"
    "fmt"
    "net/http"
    "html"
    "log"
    "os"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
    })

    serverPort := os.Getenv("SERVER_PORT")
    serverUrl := fmt.Sprintf(":%s", serverPort)
    log.Printf("Starting server on address %s", serverUrl)
    log.Fatal(http.ListenAndServe(serverUrl, nil))
}
