package models

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/url"
	"time"
)

type Visit struct {
	ID        uint      `gorm:"primarykey" json:"-"`
	UpdatedAt time.Time `json:"-"`
	CreatedAt time.Time
	UrlID     uint `json:"-"`
	IP        string
	UA        string
}

type Url struct {
	ID          uint      `gorm:"primarykey" json:"-"`
	UpdatedAt   time.Time `json:"-"`
	CreatedAt   time.Time
	Destination string
	Identifier  string `gorm:"unique;not null"`
	Label       string
	Vists       []Visit `gorm:"PRELOAD:true"`
}

const idLength = 5

func GenerateID() (string, error) {
	rnd := make([]byte, idLength)
	_, err := rand.Read(rnd)

	return hex.EncodeToString(rnd), err
}

func GetId(URL string) string {
	if len(URL) < idLength*2 {
		return ""
	}

	return URL[len(URL)-idLength*2:]
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

	newVisit := &Visit{UrlID: u.ID, IP: IP, UA: UA}

	if err = db.Create(newVisit).Error; err != nil {
		return
	}

	u.Vists = append(u.Vists, *newVisit)

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
	u.Identifier, err = GenerateID()
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

	d := db.Preload("Vists").Where("identifier = ?", Iden).First(&u)
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
