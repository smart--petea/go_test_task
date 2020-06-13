package main

import (
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"
    "github.com/joho/godotenv"
    "fmt"
    "net/http"
    "log"
    "os"
    "encoding/json"
)

func main() {
    var err error

    err = godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    db, err := connectPostgres()
    if err != nil {
        log.Fatal(err)
    }

    gormAutoMigrate(db)
    configureHttpEndpoints(db)

    err = runServer()
    if err != nil {
        log.Fatal()
    }
}

func runServer() error {
    serverPort := os.Getenv("SERVER_PORT")
    serverUrl := fmt.Sprintf("127.0.0.1:%s", serverPort)
    log.Printf("Server started at %s", serverUrl)

    return http.ListenAndServe(serverUrl, nil)
}

func configureHttpEndpoints(db *gorm.DB) {
    baseHandler := NewBaseHandler(db)

    http.HandleFunc("/category", baseHandler.getCategories)
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
    Name string `gorm:"type:varchar(100)":unique_index json:"name"`
}

type Good struct {
    gorm.Model
    Name string
    Categories []Category `gorm:"foreignkey:Id"`
}

func gormAutoMigrate(db *gorm.DB) {
    db.AutoMigrate(&Category{})
    db.AutoMigrate(&Good{})
}

type BaseHandler struct {
    db *gorm.DB
}

func NewBaseHandler(db *gorm.DB) *BaseHandler {
    return &BaseHandler{
        db: db,
    }
}

func (h *BaseHandler) getCategories(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
        case http.MethodGet:
            categories := []Category{}
            h.db.Find(&categories)

            obj, err := json.Marshal(categories)
            if err != nil {
                //add log
                http.Error(w, "error at json marshalling", 500)
                return
            } else {
                w.Write(obj)
            }

        case http.MethodPost:
            if r.Body == nil {
                http.Error(w, "Please send a request body", 400)
                log.Println("Empty body")
                return
            }

            var category Category
            err := json.NewDecoder(r.Body).Decode(&category)
            if err != nil {
                log.Println("Error: ", err)
                http.Error(w, err.Error(), 400)
                return
            }

            err = h.db.Create(&category).Error
            if err != nil {
                log.Println("Error: ", err)
                http.Error(w, err.Error(), 400)
                return
            }

            categoryJson, err := json.Marshal(category)
            if err != nil {
                log.Println("Error: ", err)
                http.Error(w, err.Error(), 400)
                return
            }

            w.Write(categoryJson)

        case http.MethodDelete:
        case http.MethodPut:
        default:

    }

}
