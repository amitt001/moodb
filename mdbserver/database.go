package server

import "github.com/amitt001/moodb/memtable"

type database struct {
	db   *memtable.KVStore
	name string
}

func (d *database) Get(key string) (string, error) {
	return d.db.Get(key)
}

func (d *database) Set(key, value string) (string, error) {
	return d.db.Create(key, value)
}

func (d *database) Del(key string) (string, error) {
	return d.db.Delete(key)
}

func newDb(name string) *database {
	return &database{db: memtable.NewDB(), name: name}
}
