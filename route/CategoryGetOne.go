package route

import (
    "net/http"
    "github.com/jinzhu/gorm"
    "log"
    "github.com/smart--petea/go_test_task/entity"
    "encoding/json"
)

type CategoryGetOne struct {}
func (*CategoryGetOne) Run(
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

    categoryJson, err := json.Marshal(category)
    if err != nil {
        log.Println("Error: ", err)
        http.Error(w, err.Error(), 400)
        return
    }
 
    w.Write(categoryJson)
}
