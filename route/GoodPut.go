package route

import (
    "net/http"
    "github.com/jinzhu/gorm"
    "log"
    "github.com/smart--petea/go_test_task/entity"
    "encoding/json"
    "fmt"
    "strconv"
)

type GoodPut struct {}

func (GoodPut) Run(
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

    if r.Body == nil {
        http.Error(w, "No request body", http.StatusBadRequest)
        log.Println("Empty body")
        return
    }

    fields := make(map[string]interface{})
    err := json.NewDecoder(r.Body).Decode(&fields)
    if err != nil {
        log.Println("Error: ", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    name, ok := fields["name"]
    if ok {
        good.Name = fmt.Sprintf("%v", name)
    } 

    price, ok := fields["price"]
    if ok {
        priceS := fmt.Sprintf("%v", price)
        priceFloat, err := strconv.ParseFloat(priceS, 32)
        if err != nil {
            log.Println("Price could not be converted to float type")
            http.Error(w, "Price could not be converted to float type", http.StatusBadRequest)
            return
        }
        good.Price = float32(priceFloat)
    }

    categoryIdsArray, ok := fields["categories"]
    if ok {
        categories := ""
        for _, categoryId := range categoryIdsArray.([]interface{}) {
            var category entity.Category
            db.Where("id = ?", categoryId).First(&category)
            if category.ID == 0 {
                msg := fmt.Sprintf("Category %v not found", categoryId)
                log.Println(msg)
                http.Error(w, msg, http.StatusBadRequest)
                return
            }

            categories = categories + "," + strconv.Itoa(int(category.ID))
        }
        good.Categories = categories[1:]
    }

    err = db.Save(&good).Error
    if err != nil {
        log.Println("Error: ", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
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
