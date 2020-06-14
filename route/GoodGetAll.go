package route

import (
    "net/http"
    "github.com/jinzhu/gorm"
    "github.com/smart--petea/go_test_task/entity"
    "encoding/json"
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
        //add log
        http.Error(w, "error at json marshalling", 500)
        return
    } else {
        w.Write(obj)
    }
}
