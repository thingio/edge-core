package storage

import (
	"bytes"
	"github.com/boltdb/bolt"
	"github.com/juju/errors"
	"github.com/thingio/edge-core/common/log"
	"github.com/thingio/edge-core/common/proto/resource"
	"os"
	"path/filepath"
	"time"
)

type BoltStorage struct {
	File       string
	Timeout	   time.Duration
	db         *bolt.DB
	bucketName []byte
}

func (s *BoltStorage) Init() error {
	log.Infof("boltdb init file: %+v", s.File)

	// create file if not exist
	if err := os.MkdirAll(filepath.Dir(s.File), 0755); err != nil {
		return err
	}

	// open db connection
	db, err := bolt.Open(s.File, 0600, &bolt.Options{Timeout: s.Timeout})
	if err != nil {
		return err
	}

	s.db = db
	s.bucketName = []byte("resource")

	// create database if not exist
	err = s.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(s.bucketName)
		return err
	})
	return err
}

func (s *BoltStorage) Put(r *resource.Resource) error {

	log.Infof("put resource: {%+v}", r)
	value, err := resource.MarshalResource(r)
	if err != nil {
		return err
	}

	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.bucketName)
		return b.Put([]byte(r.Key.String()), value)
	})

}

func (s *BoltStorage) Get(key resource.Key) (result *resource.Resource, err error) {
	log.Infof("get resource by key: %s", key)

	err = s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.bucketName)

		v := b.Get([]byte(key.String()))
		if len(v) == 0 {
			return errors.NotFoundf("value %s of key %v", s.bucketName, key.String())
		}

		r, err := resource.UnmarshalResource(key.Kind, v)
		if err != nil {
			return err
		}
		result = r
		return nil
	})

	return

}

func (s *BoltStorage) List(key resource.Key) (result []*resource.Resource, err error) {
	log.Infof("list resource by key: %s", key)

	err = s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.bucketName)

		c := b.Cursor()
		prefix := []byte(key.String())
		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {

			r, err := resource.UnmarshalResource(key.Kind, v)
			if err != nil {
				return err
			}
			result = append(result, r)
		}

		return nil
	})

	return
}
