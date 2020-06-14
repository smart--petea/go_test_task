package route

import (
    "net/http"
    "github.com/jinzhu/gorm"
    "log"
    "github.com/smart--petea/go_test_task/entity"
    "encoding/json"
)

type GoodGetOne struct {}

func (GoodGetOne) Run(
    w http.ResponseWriter,
    r *http.Request,
    params map[string]string,
    db *gorm.DB,
) {
    var good entity.Good
    db.Where("id = ?", params["id"]).First(&good)

    if good.ID == 0 {
        http.Error(w, "Good not found", http.StatusNotFound)
        log.Println("Good not found")
        return
    }

    goodJson, err := json.Marshal(good)
    if err != nil {
        log.Println("Error: ", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
 
    w.Write(goodJson)
}
