package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

var cache *Cache

func OpenDataBase(path string) (err error) {

	db, err = gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		err = db.AutoMigrate(&Url{}, &Visit{})
	}

	var Urls []Url
	result := db.Find(&Urls)
	if result.Error != nil {
		return result.Error
	}

	cache = NewCache(Urls)

	return
}
