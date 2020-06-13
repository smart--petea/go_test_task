package entity

import (
    "github.com/jinzhu/gorm"
)

type Category struct {
    gorm.Model
    Name string `gorm:"type:varchar(100)":unique_index json:"name"`
}
