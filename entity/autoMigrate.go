package entity

import (
    "github.com/jinzhu/gorm"
)

func GormAutoMigrate(db *gorm.DB) {
    db.AutoMigrate(&Category{})
    db.AutoMigrate(&Good{})
}
