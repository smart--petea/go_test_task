package main

import (
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"
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

    POSTGRES_HOST := os.Getenv("POSTGRES_HOST")
    POSTGRES_PORT := os.Getenv("POSTGRES_PORT")
    POSTGRES_USER := os.Getenv("POSTGRES_USER")
    POSTGRES_DB := os.Getenv("POSTGRES_DB")
    POSTGRES_PASSWORD := os.Getenv("POSTGRES_PASSWORD")
    POSTGRES_SSLMODE := os.Getenv("POSTGRES_SSLMODE")
    postgresUrl := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", POSTGRES_HOST, POSTGRES_PORT, POSTGRES_USER, POSTGRES_DB, POSTGRES_PASSWORD, POSTGRES_SSLMODE)
    log.Printf("Try to connect to postgres: %s", postgresUrl)
    _, err = gorm.Open("postgres", postgresUrl)
    if err != nil {
        log.Fatal(err)
    }

    http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
    })

    serverPort := os.Getenv("SERVER_PORT")
    serverUrl := fmt.Sprintf(":%s", serverPort)
    log.Printf("Starting server on address %s", serverUrl)
    log.Fatal(http.ListenAndServe(serverUrl, nil))
}
