package route

import (
    "net/http"
    "github.com/jinzhu/gorm"
)

type RoutingEntryInterface interface {
    Run(w http.ResponseWriter, r *http.Request, urlMatches map[string]string, db *gorm.DB)
}

