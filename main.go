package main

import (
    "github.com/smart--petea/go_test_task/route"
    "github.com/smart--petea/go_test_task/db"
    "github.com/smart--petea/go_test_task/app"

    _ "github.com/jinzhu/gorm/dialects/postgres"
    "github.com/joho/godotenv"
    "net/http"
    "log"
)

func main() {
    var err error

    err = godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    db, err := db.NewConnection()
    if err != nil {
        log.Fatal(err)
    }

    app := app.NewApp(db)

    app.AddRoute(`^/category$`,             http.MethodGet,  route.CategoryGetAll{})
    app.AddRoute(`^/category$`,             http.MethodPost, route.CategoryPost{})
    app.AddRoute(`^/category/(?P<id>\d+)$`, http.MethodGet,  route.CategoryGetOne{})
    app.AddRoute(`^/category/(?P<id>\d+)$`, http.MethodPut,  route.CategoryPut{})
    app.AddRoute(`^/good$`,                 http.MethodGet,  route.GoodGetAll{})
    app.AddRoute(`^/good$`,                 http.MethodPost, route.GoodPost{})
    app.AddRoute(`^/good/(?P<id>\d+)$`,     http.MethodGet,  route.GoodGetOne{})
    app.AddRoute(`^/good/(?P<id>\d+)$`,     http.MethodPut,  route.GoodPut{})

    err = app.Start()
    if err != nil {
        log.Fatal(err)
    }
}
