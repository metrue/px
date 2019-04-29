package agent

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
)

const bucketName = "processes"

type IStore interface {
	Set(k string, v interface{}) error
	Get(k string) ([]byte, error)
}

type Store struct {
	db *bolt.DB
}

func NewStore(db *bolt.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Set(k string, v interface{}) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}
		data, err := json.Marshal(v)
		if err != nil {
			return err
		}

		return b.Put([]byte(k), data)
	})
}

func (s *Store) Get(k string) ([]byte, error) {
	var v []byte
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return fmt.Errorf("bucket not existsed")
		}

		v = b.Get([]byte(k))
		return nil
	})
	return v, err
}

func (s *Store) Delete(k string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return fmt.Errorf("bucket not existsed")
		}

		return b.Delete([]byte(k))
	})
}
