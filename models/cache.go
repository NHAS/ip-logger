package models

import (
	"errors"
	"sync"
)

var ErrCacheMiss = errors.New("Entry not found in cache")
var ErrCacheTrample = errors.New("Identity already exists in cache")

type Cache struct {
	sync.RWMutex
	redirectionMap map[string]Url
}

func NewCache(urls []Url) (c *Cache) {

	c.redirectionMap = make(map[string]Url, len(urls))

	for i := range urls {
		c.redirectionMap[urls[i].Identifier] = urls[i]
	}

	return c
}

func (c *Cache) Get(Iden string) (u Url, err error) {
	c.RLock()
	defer c.RUnlock()
	if u, ok := c.redirectionMap[Iden]; ok {
		return u, nil
	}

	return u, ErrCacheMiss
}

func (c *Cache) Put(u Url) error {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.redirectionMap[u.Identifier]; ok {
		return ErrCacheTrample
	}

	c.redirectionMap[u.Identifier] = u

	return nil
}

func (c *Cache) Refresh(u Url) {
	c.Lock()
	defer c.Unlock()

	c.redirectionMap[u.Identifier] = u

}

func (c *Cache) Expire(u Url) {
	c.Lock()
	defer c.Unlock()

	delete(c.redirectionMap, u.Identifier)
}
