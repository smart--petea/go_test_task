package route

import (
    "net/http"
    "github.com/jinzhu/gorm"
    "log"
    "github.com/smart--petea/go_test_task/entity"
    "encoding/json"
)

type CategoryPut struct {}
func (*CategoryPut) Run(
    w http.ResponseWriter,
    r *http.Request,
    params map[string]string,
    db *gorm.DB,
) {
    var category entity.Category
    db.Where("id = ?", params["id"]).First(&category)

    if category.ID == 0 {
        http.Error(w, "Category not found", http.StatusNotFound)
        log.Println("Category not found")
        return
    }

    if r.Body == nil {
        http.Error(w, "No request body", http.StatusBadRequest)
        log.Println("Empty body")
        return
    }

    fields := make(map[string]string)
    err := json.NewDecoder(r.Body).Decode(&fields)
    if err != nil {
        log.Println("Error: ", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    name, ok := fields["name"]
    if ok != true {
        log.Println("There is nothing to update")
        http.Error(w, "There is nothing to update", http.StatusBadRequest)
        return
    }

    db.Model(&category).Update("name", name)

    categoryJson, err := json.Marshal(category)
    if err != nil {
        log.Println("Error: ", err)
        http.Error(w, err.Error(), 400)
        return
    }
 
    w.Write(categoryJson)
}
