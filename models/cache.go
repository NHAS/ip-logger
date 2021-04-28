package models

import (
	"errors"
	"sync"
)

var ErrCacheMiss = errors.New("Entry not found in cache")
var ErrCacheTrample = errors.New("Identity already exists in cache")

type Cache struct {
	lock           sync.RWMutex
	redirectionMap map[string]Url
}

func (c *Cache) Fuck() {
	c.lock.RLock()
	defer c.lock.RUnlock()
}

func NewCache(urls []Url) (c *Cache) {

	rMap := make(map[string]Url, len(urls))

	for i := range urls {
		rMap[urls[i].Identifier] = urls[i]
	}

	return &Cache{redirectionMap: rMap}
}

func (c *Cache) Get(Iden string) (u Url, err error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if u, ok := c.redirectionMap[Iden]; ok {
		return u, nil
	}

	return u, ErrCacheMiss
}

func (c *Cache) Put(u Url) error {

	c.lock.Lock()
	defer c.lock.Unlock()
	if _, ok := c.redirectionMap[u.Identifier]; ok {
		return ErrCacheTrample
	}

	c.redirectionMap[u.Identifier] = u

	return nil
}

func (c *Cache) Refresh(u Url) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.redirectionMap[u.Identifier] = u

}

func (c *Cache) Expire(u Url) {
	c.lock.Lock()
	defer c.lock.Unlock()

	delete(c.redirectionMap, u.Identifier)
}
