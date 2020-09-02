package database

import (
	"github.com/dgraph-io/badger"
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
	opts := badger.DefaultOptions(path)
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
			return err
		}

		err = item.Value(func(val []byte) error {
			return nil
		})
		if err != nil {
			log.Printf("ERROR: database get %s", err)
			stat.Database("get", stat.STATEFAILURE)
			return err
		}

		valCopy, err = item.ValueCopy(nil)

		return nil
	})
	if err != nil {
		log.Printf("ERROR: database get %s", err)
		stat.Database("get", stat.STATEFAILURE)
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
		return keys, err
	}

	stat.Database("get", stat.STATEOK)
	d.db.Close()
	return keys, nil
}
