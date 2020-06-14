package db

import (
    "github.com/jinzhu/gorm"
    "fmt"
    "os"
    "log"
)

func NewConnection() (*gorm.DB, error) {
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
    db, err := gorm.Open("postgres", postgresUrl)
    if err != nil {
        return nil, err
    }
    log.Printf("Connected to postgres")

    return db, nil
}
