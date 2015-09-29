package todoapp

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
)

const dbfilename = "tododb"
const bucketname = "todobucket"

type boltbackend struct {
	db *bolt.DB
}

func newBoltBackend(pdir string) (*boltbackend, error) {

	db, err := bolt.Open(filepath.Join(pdir, dbfilename), 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("error opening boltdb: %s", err)
	}

	if err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketname))
		if err != nil {
			return fmt.Errorf("error creating bucket: %s", err)
		}
		return nil
	}); err != nil {
		db.Close()
		return nil, fmt.Errorf("error opening boltdb: %s", err)
	}

	return &boltbackend{
		db: db,
	}, nil

}

func (b *boltbackend) fetchAll() ([]item, error) {
	ret := []item{}
	if err := b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketname))
		if err := b.ForEach(func(k, v []byte) error {
			var i item
			if err := gob.NewDecoder(bytes.NewReader(v)).Decode(&i); err != nil {
				return fmt.Errorf("boltbackend.fetchAll: Error decoding item %s: %s", k, err)
			}
			ret = append(ret, i)
			return nil
		}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return ret, nil
}

func (b *boltbackend) create(i item) error {
	if err := b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketname))
		buf := bytes.Buffer{}
		if err := gob.NewEncoder(&buf).Encode(i); err != nil {
			return fmt.Errorf("error encoding item: %s", err)
		}
		if err := b.Put([]byte(i.ID), buf.Bytes()); err != nil {
			return fmt.Errorf("error putting item %v: %s", i, err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("boltbackend.create: error creating item: %s", err)
	}
	return nil
}

func (b *boltbackend) delete(id string) error {
	if err := b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketname))
		if err := b.Delete([]byte(id)); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return fmt.Errorf("boltbackend.delete: error deleting item: %s", err)
	}
	return nil
}

func (b *boltbackend) stop() {
	b.db.Close()
}
