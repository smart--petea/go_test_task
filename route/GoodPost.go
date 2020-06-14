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

type GoodPost struct {}

func (GoodPost) Run(
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

    var good entity.Good
    obj := make(map[string]interface{})
    err := json.NewDecoder(r.Body).Decode(&obj)
    if err != nil {
        log.Println("Error: ", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    name, ok := obj["name"]
    if !ok {
        log.Println("Name not provided")
        http.Error(w, "Name not provided", http.StatusBadRequest)
        return
    }

    nameS, ok := name.(string)
    if !ok {
        log.Println("Name is not of string type")
        http.Error(w, "Name is not of string type", http.StatusBadRequest)
        return
    }

    good.Name = nameS

    price, ok := obj["price"]
    if !ok {
        log.Println("Price not provided")
        http.Error(w, "Price not provided", http.StatusBadRequest)
        return
    }

    priceS := fmt.Sprintf("%v", price)
    priceFloat, err := strconv.ParseFloat(priceS, 32)
    if err != nil {
        log.Println("Price could not be converted to float type")
        http.Error(w, "Price could not be converted to float type", http.StatusBadRequest)
        return
    }

    good.Price = float32(priceFloat)

    categoryIdsArray, ok := obj["categories"]
    if !ok {
        log.Println("Categories not provided")
        http.Error(w, "Categories not provided", http.StatusBadRequest)
        return
    }


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

    log.Println(good)
    err = db.Create(&good).Error
    if err != nil {
        log.Println("Error: ", err)
        http.Error(w, err.Error(), 400)
        return
    }

    goodJson, err := json.Marshal(good)
    if err != nil {
        log.Println("Error: ", err)
        http.Error(w, err.Error(), 400)
        return
    }

    w.Write(goodJson)
}
