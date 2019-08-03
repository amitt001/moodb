package server

import (
	"log"
	"github.com/amitt001/moodb/wal"
	"github.com/amitt001/moodb/memtable"
)

// TODO: check can both database and KVStore use an interface?

// database: main db object. It wraps the `memtable` methods.
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
	db := &database{db: memtable.NewDB(), name: name}
	log.Println("Starting DB recovery")
	func() {
		w, err := wal.InitWal("/Users/amittripathi/codes/go/src/github.com/amitt001/moodb/data", true)

		if err != nil {
			if err == wal.ErrWalNotFound {
				return
			}
			log.Fatal(err)
		}
		rChan, err := w.Replay()
		if err != nil {
			log.Fatal(err)
		}
		for record := range rChan {
			db.Set(record.Data.GetKey(), record.Data.GetValue())
		}
	}()
	log.Println("DB recovery finished")
	return db
}
