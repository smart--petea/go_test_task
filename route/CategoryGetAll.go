package route

import (
    "net/http"
    "github.com/jinzhu/gorm"
    "github.com/smart--petea/go_test_task/entity"
    "encoding/json"
    "log"
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
        log.Print(err)
        http.Error(w, "error at json marshalling", http.StatusInternalServerError)
        return
    } else {
        w.Write(obj)
    }
}
