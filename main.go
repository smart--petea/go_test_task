package main

import (
    "github.com/smart--petea/go_test_task/route"
    "github.com/smart--petea/go_test_task/entity"

    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"
    "github.com/joho/godotenv"
    "fmt"
    "net/http"
    "log"
    "os"
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

    entity.GormAutoMigrate(db)

    err = runServer(db)
    if err != nil {
        log.Fatal()
    }
}

func runServer(db *gorm.DB) error {
    serverPort := os.Getenv("SERVER_PORT")
    serverUrl := fmt.Sprintf("127.0.0.1:%s", serverPort)
    log.Printf("Server started at %s", serverUrl)

    routing := make(map[string]map[*regexp.Regexp] route.RoutingEntryInterface)
    routing[http.MethodGet] = make(map[*regexp.Regexp] route.RoutingEntryInterface)
    routing[http.MethodPut] = make(map[*regexp.Regexp] route.RoutingEntryInterface)
    routing[http.MethodPost] = make(map[*regexp.Regexp] route.RoutingEntryInterface)

    categoryRegex := regexp.MustCompile(`^/category$`)
    routing[http.MethodGet][categoryRegex] = &route.CategoryGetAll{}
    routing[http.MethodPost][categoryRegex] = &route.CategoryPost{}

    categoryIdRegex := regexp.MustCompile(`^/category/(?P<id>\d+)$`)
    routing[http.MethodGet][categoryIdRegex] = &route.CategoryGetOne{}
    routing[http.MethodPut][categoryIdRegex] = &route.CategoryPut{}

    goodRegex := regexp.MustCompile(`^/good$`)
    routing[http.MethodGet][goodRegex] = &route.GoodGetAll{}
    routing[http.MethodPost][goodRegex] = &route.GoodPost{}

    goodIdRegex := regexp.MustCompile(`^/good/(?P<id>\d+)$`)
    routing[http.MethodGet][goodIdRegex] = &route.GoodGetOne{}
    routing[http.MethodPut][goodIdRegex] = &route.GoodPut{}

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

            routingEntry.Run(w, r, result, db)
            return
        }

        http.Error(w, "no suitable route found", 400)
    })

    return http.ListenAndServe(serverUrl, nil)
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

