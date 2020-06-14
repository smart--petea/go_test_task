package route

import (
    "net/http"
    "github.com/jinzhu/gorm"
)

type RouteInterface interface {
    Run(
        w http.ResponseWriter,
        r *http.Request,
        urlMatches map[string]string,
        db *gorm.DB,
    )
}

