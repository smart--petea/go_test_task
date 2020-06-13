package entity

import (
    "github.com/jinzhu/gorm"
)

type Good struct {
    gorm.Model
    Name string
    Categories []Category `gorm:"foreignkey:Id"`
}
