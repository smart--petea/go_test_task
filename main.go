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
    "regexp"
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

    err = runServer(db)
    if err != nil {
        log.Fatal()
    }
}

type RoutingEntryInterface interface {
    run(w http.ResponseWriter, r *http.Request, urlMatches map[string]string, db *gorm.DB)
}

func runServer(db *gorm.DB) error {
    serverPort := os.Getenv("SERVER_PORT")
    serverUrl := fmt.Sprintf("127.0.0.1:%s", serverPort)
    log.Printf("Server started at %s", serverUrl)

    routing := make(map[string]map[*regexp.Regexp]RoutingEntryInterface)
    routing[http.MethodGet] = make(map[*regexp.Regexp]RoutingEntryInterface)
    routing[http.MethodPut] = make(map[*regexp.Regexp]RoutingEntryInterface)
    routing[http.MethodPost] = make(map[*regexp.Regexp]RoutingEntryInterface)
    routing[http.MethodDelete] = make(map[*regexp.Regexp]RoutingEntryInterface)

    categoryRegex := regexp.MustCompile(`^/category$`)
    routing[http.MethodGet][categoryRegex] = &CategoryGetAll{}
    routing[http.MethodPost][categoryRegex] = &CategoryPost{}

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s", r.Method, r.URL.Path)
        _, ok := routing[r.Method]
        if !ok {
            http.Error(w, "Not a supported http method", 500)
            return 
        }

        for routingRegex, routingEntry := range routing[r.Method] {
            match := routingRegex.FindStringSubmatch(r.URL.Path)
            if len(match) == 0 {
                continue
            }

            result := make(map[string]string)
            for i, name := range routingRegex.SubexpNames() {
                if i != 0 && name != "" {
                    result[name] = match[i]
                }
            }

            routingEntry.run(w, r, result, db)
            return
        }

        http.Error(w, "no suitable route found", 400)
    })

    return http.ListenAndServe(serverUrl, nil)
}

/*
func configureHttpEndpoints(db *gorm.DB) {
    baseHandler := NewBaseHandler(db)

    http.HandleFunc("/", baseHandler.getCategories)
}
*/

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

type CategoryGetAll struct {}
func (*CategoryGetAll) run(
    w http.ResponseWriter,
    r *http.Request,
    urlMatches map[string]string,
    db *gorm.DB,
) {
    categories := []Category{}
    db.Find(&categories)

    obj, err := json.Marshal(categories)
    if err != nil {
        //add log
        http.Error(w, "error at json marshalling", 500)
        return
    } else {
        w.Write(obj)
    }
}

type CategoryPost struct {}
func (*CategoryPost) run(
    w http.ResponseWriter,
    r *http.Request,
    urlMatches map[string]string,
    db *gorm.DB,
) {
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

    err = db.Create(&category).Error
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
}
