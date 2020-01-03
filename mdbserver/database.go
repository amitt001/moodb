package server

import (
	"github.com/amitt001/moodb/mdbserver/mdbserverpb"
	"github.com/amitt001/moodb/memtable"
	"github.com/amitt001/moodb/wal"
	"github.com/golang/protobuf/proto"
	"log"
	"sync"
)

const (
	// Maybe this should be an ENUM
	active   = "ACTIVE"
	recovery = "RECOVERY"
)

// TODO: check can both database and KVStore use an interface?

// database: main db object. It wraps the `memtable` methods.
type database struct {
	db      *memtable.KVStore
	name    string
	mode    string // mode DB is currently running in. Can be recovery/active
	mu      sync.Mutex
	rWalObj *wal.Wal // used during recovery
	walObj  *wal.Wal
}

func (d *database) logRecord(cmd, key, value string) error {
	record, err := proto.Marshal(&mdbserverpb.Record{Cmd: cmd, Key: key, Value: value})
	if err != nil {
		return err
	}
	err = d.walObj.Write(record)
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
	err := d.logRecord("DEL", key, "")
	if err != nil {
		log.Fatal(err)
	}
	return d.db.Delete(key)
}

func (d *database) recoverySet(key, value string) (string, error) {
	return d.db.Create(key, value)
}

// toggleMode changes the db mode from active to recovery if in active mode
func (d *database) setMode(mode string) {
	d.mode = mode
}

func newDb(name, walDir string) *database {
	db := &database{db: memtable.NewDB(), name: name}
	db.mu.Lock()
	defer db.mu.Unlock()
	log.Println("Starting DB recovery")
	var err error
	func() {
		db.setMode(recovery)
		defer db.setMode(active)
		recover := true
		// Open recovery WAL file
		db.rWalObj, err = wal.Open(walDir)
		if err != nil {
			if err == wal.ErrWalNotFound {
				recover = false
			} else {
				log.Fatalf("Recovery: %s", err)
			}
		}
		// Open new WAL tmp file
		db.walObj, err = wal.New(walDir)
		if err != nil {
			log.Fatalf("Recovery: %s", err)
		}
		if recover {
			rChan := db.rWalObj.Read()
			for record := range rChan {
				recordData := &mdbserverpb.Record{}
				err = proto.Unmarshal(record.Data, recordData)
				// TODO handle error here
				switch recordData.Cmd {
				case "SET":
					db.Set(recordData.GetKey(), recordData.GetValue())
				case "DEL":
					db.Del(recordData.GetKey())
				default:
					log.Fatal("Recovery: Invalid command")
				}
			}
			// Free the recovery WAL object
			db.rWalObj.Close()
			db.rWalObj = nil
		}
	}()

	log.Println("DB recovery finished")
	return db
}
