package entity

import (
    "github.com/jinzhu/gorm"
)

type Good struct {
    gorm.Model
    Name string `gorm:"type:varchar(100);unique_index"json:"name"`
    Price float32 `gorm:"type:decimal(100,2)" json:"price"`
    Categories string `gorm:"type:varchar(100);" json:"categories"`
}
