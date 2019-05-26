package database

import (
	"errors"

	"go.etcd.io/bbolt"
	"go.ybk.im/homepage/internal/app/types"
)

var siteBucket = []byte("site")

func PutSite(site *types.Site) error {
	db, err := instance()
	if err != nil {
		return err
	}

	return db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(siteBucket)
		if err != nil {
			return err
		}

		err = b.Put([]byte("title"), []byte(site.Title))
		if err != nil {
			return err
		}

		err = b.Put([]byte("summary"), []byte(site.Summary))
		if err != nil {
			return err
		}

		return nil
	})
}

func GetSite(site *types.Site) error {
	if site == nil {
		return errors.New("site object is nil")
	}

	db, err := instance()
	if err != nil {
		return err
	}

	return db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(siteBucket)
		if b == nil {
			return nil
		}

		title := b.Get([]byte("title"))
		if title != nil {
			site.Title = string(title)
		}

		summary := b.Get([]byte("summary"))
		if summary != nil {
			site.Summary = string(summary)
		}

		return nil
	})
}
