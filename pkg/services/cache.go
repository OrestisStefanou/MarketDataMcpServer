package services

import (
	"encoding/json"
	"time"

	badger "github.com/dgraph-io/badger/v4"
)

type CacheService interface {
	Get(key string, target interface{}) error
	Set(key string, value interface{}, ttl time.Duration) error
	Delete(key string) error
}

type BadgerCacheService struct {
	db *badger.DB
}

func NewBadgerCacheService() (*BadgerCacheService, error) {
	db, err := badger.Open(badger.DefaultOptions("cache.db"))
	if err != nil {
		return nil, err
	}
	return &BadgerCacheService{db: db}, nil
}

func (c *BadgerCacheService) Get(key string, target interface{}) error {
	var data []byte
	err := c.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		data, err = item.ValueCopy(nil)
		return err
	})
	if err != nil {
		return err
	}

	return json.Unmarshal(data, target)
}

func (c *BadgerCacheService) Set(key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = c.db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(key), data).WithTTL(ttl)
		return txn.SetEntry(e)
	})

	return err
}

func (c *BadgerCacheService) Delete(key string) error {
	return c.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}
