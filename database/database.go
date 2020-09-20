package database

import (
	"github.com/dgraph-io/badger"
	"github.com/rebelit/gome-schedule/common/config"
	"github.com/rebelit/gome-schedule/common/stat"
	"log"
)

var (
	dataDirectory = "/tmp/badger"
)

type Database interface {
	Set(key string, value []byte) error
	Delete(key string) error
	DeleteAll() error
	Get(key string) ([]byte, error)
	GetAllKeys() ([]string, error)
	Close() error
}

type Badger struct {
	db *badger.DB
}

func Open(path string) (Database, error) {
	if path == "" {
		path = dataDirectory
	}

	d := Badger{}
	opts := badger.Options{}
	if config.App.FullMemory {
		opts = badger.DefaultOptions(path)
	} else {
		opts = badger.Options{
			Dir:                     path,
			ValueDir:                path,
			SyncWrites:              true,
			TableLoadingMode:        0,
			ValueLogLoadingMode:     0,
			LevelOneSize:            256 << 20,
			LevelSizeMultiplier:     10,
			NumVersionsToKeep:       1,
			ReadOnly:                false,
			Truncate:                false,
			Logger:                  nil,
			EventLogging:            true,
			MaxTableSize:            64 << 20,
			MaxLevels:               7,
			ValueThreshold:          32,
			NumMemtables:            5,
			NumLevelZeroTables:      5,
			NumLevelZeroTablesStall: 10,
			ValueLogFileSize:        1<<30 - 1,
			ValueLogMaxEntries:      1000,
			NumCompactors:           2,
			CompactL0OnClose:        true,
			LogRotatesToFlush:       2,
			VerifyValueChecksum:     false,
			BypassLockGuard:         false,
		}
	}

	db, err := badger.Open(opts)
	if err != nil {
		log.Printf("ERROR: database %s open %s", path, err)
		stat.Database("open", stat.STATEFAILURE)
		return d, err
	}

	stat.Database("open", stat.STATEOK)
	d.db = db
	return d, nil
}

func (d Badger) Close() error {
	err := d.db.Close()
	if err != nil {
		return err
	}

	return nil
}

func (d Badger) Set(key string, value []byte) error {
	err := d.db.Update(func(txn *badger.Txn) error {
		txn.Set([]byte(key), value)
		return nil
	})
	if err != nil {
		log.Printf("ERROR: database set %s", err)
		stat.Database("set", stat.STATEFAILURE)
		d.db.Close()
		return err
	}

	stat.Database("set", stat.STATEOK)
	d.db.Close()
	return nil
}

func (d Badger) Delete(key string) error {
	err := d.db.Update(func(txn *badger.Txn) error {
		txn.Delete([]byte(key))
		return nil
	})
	if err != nil {
		log.Printf("ERROR: database delete %s", err)
		stat.Database("delete", stat.STATEFAILURE)
		d.db.Close()
		return err
	}

	stat.Database("delete", stat.STATEOK)
	d.db.Close()
	return nil
}

func (d Badger) DeleteAll() error {
	err := d.db.DropAll()
	if err != nil {
		log.Printf("ERROR: database delete_all %s", err)
		stat.Database("delete", stat.STATEFAILURE)
		d.db.Close()
		return err
	}

	stat.Database("delete", stat.STATEOK)
	d.db.Close()
	return nil
}

func (d Badger) Get(key string) ([]byte, error) {
	var valCopy []byte
	err := d.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			log.Printf("ERROR: database get %s", err)
			stat.Database("get", stat.STATEFAILURE)
			d.db.Close()
			return err
		}

		err = item.Value(func(val []byte) error {
			return nil
		})
		if err != nil {
			log.Printf("ERROR: database get %s", err)
			stat.Database("get", stat.STATEFAILURE)
			d.db.Close()
			return err
		}

		valCopy, err = item.ValueCopy(nil)
		return nil
	})
	if err != nil {
		log.Printf("ERROR: database get %s", err)
		stat.Database("get", stat.STATEFAILURE)
		d.db.Close()
		return nil, err
	}

	stat.Database("get", stat.STATEOK)
	d.db.Close()
	return valCopy, nil
}

func (d Badger) GetAllKeys() ([]string, error) {
	keys := []string{}
	err := d.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			keys = append(keys, string(k))
		}
		return nil
	})
	if err != nil {
		log.Printf("ERROR: database get_all %s", err)
		stat.Database("get", stat.STATEFAILURE)
		d.db.Close()
		return keys, err
	}

	stat.Database("get", stat.STATEOK)
	d.db.Close()
	return keys, nil
}
