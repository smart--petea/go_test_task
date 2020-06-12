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

    db, err := connectPostgres()
    if err != nil {
        log.Fatal(err)
    }

    gormAutoMigrate(db)

    http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
    })

    serverPort := os.Getenv("SERVER_PORT")
    serverUrl := fmt.Sprintf("127.0.0.1:%s", serverPort)
    log.Printf("Server started at %s", serverUrl)
    log.Fatal(http.ListenAndServe(serverUrl, nil))
}

func connectPostgres() (db *gorm.DB, err error) {
    postgresUrl := fmt.Sprintf(
        "host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
        os.Getenv("POSTGRES_HOST"),
        os.Getenv("POSTGRES_PORT"),
        os.Getenv("POSTGRES_USER"),
        os.Getenv("POSTGRES_DB"),
        os.Getenv("POSTGRES_PASSWORD"),
        os.Getenv("POSTGRES_SSLMODE"),
    )

    log.Printf("Connecting to postgres...")
    db, err = gorm.Open("postgres", postgresUrl)
    if err != nil {
        return nil, err
    }
    log.Printf("Connected to postgres")

    return db, nil
}

type Category struct {
    gorm.Model
    Name string
}

func gormAutoMigrate(db *gorm.DB) {
    db.AutoMigrate(&Category{})
}
