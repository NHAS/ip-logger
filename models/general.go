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
		return

	}

	err = db.AutoMigrate(&Url{}, &Visit{})
	if err != nil {
		return
	}

	var Urls []Url
	if err = db.Find(&Urls).Error; err != nil {
		return
	}

	cache = NewCache(Urls)

	return
}
