package server

import (
	"errors"
	"github.com/amitt001/moodb/memtable"
	"github.com/amitt001/moodb/wal"
	"github.com/amitt001/moodb/wal/walpb"
	"log"
	"sync"
)

const (
	// Maybe this should be an ENUM
	active   = "ACTIVE"
	recovery = "RECOVERY"
	walFile  = "/Users/amittripathi/codes/go/src/github.com/amitt001/moodb/data"
)

var (
	ErrWalNotFound = errors.New("Wal not present")
)

// TODO: check can both database and KVStore use an interface?

// database: main db object. It wraps the `memtable` methods.
type database struct {
	db      *memtable.KVStore
	name    string
	mode    string // mode DB is currently running in. Can be recovery/active
	mu      sync.Mutex
	rwalObj *wal.Wal // used during recovery
	walObj  *wal.Wal
}

func (d *database) logRecord(cmd, key, value string) error {
	record := &walpb.Data{Cmd: cmd, Key: key, Value: value}
	err := d.walObj.AppendLog(record)
	return err
}

func (d *database) Get(key string) (string, error) {
	return d.db.Get(key)
}

func (d *database) Set(key, value string) (string, error) {
	err := d.logRecord("SET", key, value)
	if err != nil {
		log.Fatal(err)
	}
	return d.db.Create(key, value)
}

func (d *database) Del(key string) (string, error) {
	return d.db.Delete(key)
}

func (d *database) recoverySet(key, value string) (string, error) {
	return d.db.Create(key, value)
}

// toggleMode changes the db mode from active to recovery if in active mode
func (d *database) setMode(mode string) {
	d.mode = mode
}

func (d *database) initWal(recovery bool) error {
	w, err := wal.InitWal(walFile, recovery)
	if w.IsWalPresent() == false {
		return ErrWalNotFound
	}
	d.walObj = w
	return err
}

func newDb(name string) *database {
	db := &database{db: memtable.NewDB(), name: name}
	db.mu.Lock()
	defer db.mu.Unlock()
	log.Println("Starting DB recovery")
	func() {
		db.setMode(recovery)
		defer db.setMode(active)
		err := db.initWal(true)
		if err != nil {
			if err == ErrWalNotFound {
				return
			}
			log.Fatal(err)
		}
		rChan, err := db.walObj.Replay()
		if err != nil {
			log.Fatal(err)
		}
		for record := range rChan {
			k := record.Data.GetKey()
			log.Println(k)
			db.recoverySet(k, record.Data.GetValue())
		}
	}()

	err := db.initWal(false)
	if err != nil {
		if err != ErrWalNotFound {
			log.Fatal(err)
		}
	}

	log.Println("DB recovery finished")
	return db
}
