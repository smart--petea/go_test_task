package route

import (
    "net/http"
    "github.com/jinzhu/gorm"
    "github.com/smart--petea/go_test_task/entity"
    "encoding/json"
)

type CategoryGetAll struct {}

func (CategoryGetAll) Run(
    w http.ResponseWriter,
    r *http.Request,
    params map[string]string,
    db *gorm.DB,
) {
    categories := []entity.Category{}
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
