package bolt

import (
	"fmt"

	bolt "go.etcd.io/bbolt"
)

const (
	ShareTokenBucketName = "share_token"
)

type KeyValueStore struct {
	DB *bolt.DB
}

func NewConn() (*KeyValueStore, error) {
	db, err := bolt.Open("/Users/facai/Documents/code/apicat/apicat_open/data/apicat_bolt.db", 0600, nil)
	if err != nil {
		return nil, err
	}

	return &KeyValueStore{DB: db}, nil
}

func (kvs *KeyValueStore) Close() error {
	return kvs.DB.Close()
}

func (kvs *KeyValueStore) Get(bucketName, key []byte) ([]byte, error) {
	defer kvs.Close()

	var value []byte
	err := kvs.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}

		value = bucket.Get(key)
		return nil
	})

	if err != nil {
		return nil, err
	}
	return value, nil
}

func (kvs *KeyValueStore) Put(bucketName, key, value []byte) error {
	defer kvs.Close()

	err := kvs.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			return err
		}

		err = bucket.Put(key, value)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}
