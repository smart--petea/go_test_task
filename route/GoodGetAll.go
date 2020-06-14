package route

import (
    "net/http"
    "github.com/jinzhu/gorm"
    "github.com/smart--petea/go_test_task/entity"
    "encoding/json"
    "log"
)

type GoodGetAll struct {}

func (GoodGetAll) Run(
    w http.ResponseWriter,
    r *http.Request,
    params map[string]string,
    db *gorm.DB,
) {
    goods := []entity.Good{}
    db.Find(&goods)

    obj, err := json.Marshal(goods)
    if err != nil {
        log.Println(err)
        http.Error(w, "error at json marshalling", http.StatusInternalServerError)
        return
    } else {
        w.Write(obj)
    }
}
