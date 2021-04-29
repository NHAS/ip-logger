package models

import (
	"fmt"
	"net/url"
	"time"

	"github.com/NHAS/ip-logger/util"
)

type Visit struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UrlID     uint
	IP        string
	UA        string
}

type Url struct {
	ID          uint `gorm:"primarykey"`
	CreatedAt   time.Time
	Destination string
	Identifier  string `gorm:"unique;not null"`
	Label       string
	Vists       []Visit `gorm:"PRELOAD:true"`
}

func NewVisit(Identifier, IP, UA string) (err error) {
	u, err := cache.Get(Identifier)
	if err != nil {
		if err != ErrCacheMiss {
			return
		}

		if err = db.Where("identifier = ?", Identifier).First(&u).Error; err != nil {
			return
		}
	}

	if err = db.Create(&Visit{UrlID: u.ID, IP: IP, UA: UA}).Error; err != nil {
		return
	}

	cache.Refresh(u)

	return
}

func NewUrl(Dest, Label string) (id string, err error) {
	_, err = url.Parse(Dest)
	if err != nil {
		return
	}

	var u Url
	u.Destination = Dest
	u.Label = Label
	u.Identifier, err = util.GenerateID()
	if err != nil {
		return
	}

	if err = db.Create(&u).Error; err != nil {
		return
	}

	return u.Identifier, cache.Put(u)
}

func GetUrl(Iden string) (u Url, err error) {
	if u, err = cache.Get(Iden); err == nil {
		return
	}

	d := db.Where("identifier = ?", Iden).First(&u)
	if d.Error != nil {
		return u, d.Error
	}

	err = cache.Put(u)

	return
}

func GetAllUrls() (us []Url, err error) {
	if err = db.Preload("Vists").Find(&us).Error; err != nil {
		return us, err
	}

	return
}

func DeleteUrl(Ident string) (err error) {
	var u Url
	tx := db.Where("identifier = ?", Ident).Find(&u)
	if err = tx.Error; err != nil {
		return tx.Error
	}

	tx = db.Delete(&u)
	if err = tx.Error; err != nil {
		return
	}

	if tx.RowsAffected == 0 {
		return fmt.Errorf("Unknown entry '%s'", Ident)
	}

	cache.Expire(u)

	return db.Delete(&Visit{}, "url_id = ?", u.ID).Error
}
