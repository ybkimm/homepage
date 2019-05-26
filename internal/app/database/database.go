package database

import (
	"errors"

	"go.etcd.io/bbolt"
)

var dbInst *bbolt.DB

func Init(fileName string) error {
	db, err := bbolt.Open(fileName, 0666, nil)
	if err != nil {
		return err
	}

	dbInst = db
	return nil
}

func Close() {
	if dbInst == nil {
		return
	}
	_ = dbInst.Close()
}

func instance() (*bbolt.DB, error) {
	if dbInst == nil {
		return nil, errors.New("database is not initialized")
	}
	return dbInst, nil
}
