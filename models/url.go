package models

import (
	"net/url"

	"github.com/NHAS/ip-logger/util"
	"gorm.io/gorm"
)

type Visit struct {
	gorm.Model
	UrlID uint
	IP    string
}

type Url struct {
	gorm.Model
	Destination string
	Identifier  string `gorm:"unique;not null"`
	Label       string
	Vists       []Visit `gorm:"PRELOAD:true"`
}

func NewVisit(Identifier, IP string) (err error) {
	u, err := cache.Get(Identifier)
	if err != nil {
		if err := db.Where("identifier = ?", Identifier).First(&u).Error; err != nil {
			return err
		}
	}

	if err := db.Create(&Visit{UrlID: u.ID, IP: IP}).Error; err != nil {
		return err
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
	if err = db.Find(&us).Error; err != nil {
		return us, err
	}

	return
}

func DeleteUrl(Ident string) (err error) {
	var u Url
	if err = db.Where("identifier = ?", Ident).Delete(&u).Error; err != nil {
		return
	}

	cache.Expire(u)

	return db.Delete(&Visit{}, "url_id = ?", u.ID).Error
}
