package entity

import (
    "github.com/jinzhu/gorm"
)

type Good struct {
    gorm.Model
    Name string `gorm:"type:varchar(100);unique_index"json:"name"`
    Price float32 `gorm:"type:decimal(100,2)":unique_index json:"name"`
    Categories []Category `gorm:"foreignkey:Id" json:"categories"`
}
