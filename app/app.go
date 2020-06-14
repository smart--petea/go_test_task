package app

import (
    "github.com/jinzhu/gorm"
    "github.com/smart--petea/go_test_task/route"
    "github.com/smart--petea/go_test_task/entity"

    "os"
    "fmt"
    "log"
    "regexp"
    "net/http"
)

type App struct {
    db *gorm.DB
    routeTable map[string]map[*regexp.Regexp] *route.RouteInterface
}

func NewApp(db *gorm.DB) App {
    routeTable := make(map[string]map[*regexp.Regexp]*route.RouteInterface)
    routeTable[http.MethodGet] = make(map[*regexp.Regexp]*route.RouteInterface)
    routeTable[http.MethodPut] = make(map[*regexp.Regexp]*route.RouteInterface)
    routeTable[http.MethodPost] = make(map[*regexp.Regexp]*route.RouteInterface)

    app := App{
        db: db,
        routeTable: routeTable,
    }

    entity.GormAutoMigrate(app.db)

    return app
}

func (app App) Start() error {
    err := app.runServer()

    return err
}

func (app *App) runServer() error {
    serverPort := os.Getenv("SERVER_PORT")
    serverUrl := fmt.Sprintf("127.0.0.1:%s", serverPort)
    log.Printf("Server started at %s", serverUrl)

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s", r.Method, r.URL.Path)
        _, ok := app.routeTable[r.Method]
        if !ok {
            http.Error(w, "Not a supported http method", 500)
            return 
        }

        for routingRegex, routingEntry := range app.routeTable[r.Method] {
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

            (*routingEntry).Run(w, r, result, app.db)
            return
        }

        http.Error(w, "no suitable route found", 400)
    })

    return http.ListenAndServe(serverUrl, nil)
}

func (app *App) AddRoute(
    urlRegex string,
    httpMethod string,
    routeObj route.RouteInterface,
) {
    app.routeTable[httpMethod][regexp.MustCompile(urlRegex)] = &routeObj
}
