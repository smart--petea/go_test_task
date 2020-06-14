package route

import (
    "net/http"
    "github.com/jinzhu/gorm"
    "log"
    "github.com/smart--petea/go_test_task/entity"
    "encoding/json"
)

type CategoryPost struct {}

func (CategoryPost) Run(
    w http.ResponseWriter,
    r *http.Request,
    params map[string]string,
    db *gorm.DB,
) {
    if r.Body == nil {
        http.Error(w, "Please send a request body", 400)
        log.Println("Empty body")
        return
    }

    var category entity.Category
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
