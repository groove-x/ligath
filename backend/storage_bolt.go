package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"go.etcd.io/bbolt"
)

type BoltStorage struct {
	db *bbolt.DB
}

func NewBoltStorage(dbPath string) (*BoltStorage, error) {
	db, err := bbolt.Open(dbPath, 0666, nil)
	if err != nil {
		return nil, err
	}

	b := &BoltStorage{
		db: db,
	}

	return b, nil
}

func (s *BoltStorage) Setup() error {
	jsons, err := filepath.Glob("./*.json")
	if err != nil {
		return fmt.Errorf("failed to setup boltdb: %v", err)
	}

	err = s.Migrate(jsons)
	if err != nil {
		return fmt.Errorf("failed to migrate JSONs: %v", err)
	}

	return nil
}

func (s *BoltStorage) Close() {
	s.db.Close()
}

func (s *BoltStorage) Migrate(jsons []string) error {
	for _, j := range jsons {
		f, err := os.Open(j)
		if err != nil {
			err = fmt.Errorf("failed to migrate: failed to open input JSON: %v", err)
			log.Print(err)
			return err
		}

		jsonBuf, err := ioutil.ReadAll(f)
		if err != nil {
			return fmt.Errorf("failed to read %s: %v", j, err)
		}
		f.Close()

		migrated, err := s.migrate(j, jsonBuf)
		if err != nil {
			return err
		}

		if migrated {
			log.Printf("migration: %s has just migrated", filepath.Base(j))
		} else {
			log.Printf("migration: %s is already migrated, skipping", filepath.Base(j))
		}
	}
	return nil
}

func (s *BoltStorage) migrate(fn string, buf []byte) (bool, error) {
	fnBase := filepath.Base(fn)
	migrated := s.getMigratedJSONs()
	for _, m := range migrated {
		if m == fnBase {
			return false, nil
		}
	}

	aptOutput := AptOutput{}
	err := json.Unmarshal(buf, &aptOutput)
	if err != nil {
		err = fmt.Errorf("failed to migrate: failed to parse %s: %v", fn, err)
		log.Print(err)
		return false, err
	}

	err = s.db.Update(func(tx *bbolt.Tx) error {
		migration, err := tx.CreateBucketIfNotExists([]byte("migration"))
		if err != nil {
			err = fmt.Errorf("failed to migrate: failed to create bucket: %v", err)
			log.Print(err)
			return err
		}

		err = migration.Put([]byte(fnBase), []byte(fnBase))
		if err != nil {
			err = fmt.Errorf("failed to migrate: failed to Put in migration bucket: %v", err)
			log.Print(err)
			return err
		}

		parsed, err := tx.CreateBucketIfNotExists([]byte(fmt.Sprintf("%s_parsed", fnBase)))
		if err != nil {
			err = fmt.Errorf("failed to migrate: failed to create bucket: %v", err)
			log.Print(err)
			return err
		}

		for _, pkg := range aptOutput.Parsed {
			marshaled, err := json.Marshal(pkg)
			if err != nil {
				err = fmt.Errorf(
					"failed to migrate: failed to marshal package %s/%s: %v",
					pkg.Name,
					pkg.Version,
					err,
				)
				log.Print(err)
				return err
			}
			parsed.Put([]byte(pkg.Name+"/"+pkg.Version), marshaled)
		}

		notParsed, err := tx.CreateBucketIfNotExists([]byte(fn + "_notparsed"))
		if err != nil {
			err = fmt.Errorf("failed to migrate: failed to create bucket: %v", err)
			log.Print(err)
			return err
		}

		for _, pkg := range aptOutput.NotParsed {
			marshaled, err := json.Marshal(pkg)
			if err != nil {
				err = fmt.Errorf(
					"failed to migrate: failed to marshal package %s/%s: %v",
					pkg.Name,
					pkg.Version,
					err,
				)
				log.Print(err)
				return err
			}
			notParsed.Put([]byte(pkg.Name+"/"+pkg.Version), marshaled)
		}
		return nil
	})

	if err != nil {
		err = fmt.Errorf("failed to migrate: %v", err)
		log.Print(err)
		return false, err
	}

	return true, nil
}

func (s *BoltStorage) getMigratedJSONs() []string {
	var l []string

	s.db.View(func(tx *bbolt.Tx) error {
		return tx.ForEach(func(name []byte, bucket *bbolt.Bucket) error {
			if string(name) == "migration" {
				bucket.ForEach(func(k, _ []byte) error {
					l = append(l, string(k))
					return nil
				})
			}
			return nil
		})
	})
	return l
}

func (s *BoltStorage) GetPackage(pkg, ver string) (*Package, error) {
	var item *Package
	noProblemo := errors.New("")

	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.ForEach(func(name []byte, bucket *bbolt.Bucket) error {
			res := bucket.Get([]byte(pkg + "/" + ver))
			if res == nil {
				return nil
			}

			item = &Package{}
			err := json.Unmarshal(res, item)
			if err != nil {
				return fmt.Errorf("found item, but failed to unmarshal: %v", err)
			}
			return noProblemo
		})
	})

	if err != nil && err.Error() != noProblemo.Error() {
		return nil, fmt.Errorf("failed to get package: %v", err)
	}

	if item == nil {
		return nil, nil
	}

	return item, nil
}

func (s *BoltStorage) getPackages(kind string) []PackageListItem {
	l := make([]PackageListItem, 0, 100)
	s.db.View(func(tx *bbolt.Tx) error {
		tx.ForEach(func(name []byte, b *bbolt.Bucket) error {
			if !bytes.HasSuffix(name, []byte(kind)) {
				return nil
			}

			b.ForEach(func(k, v []byte) error {
				sp := strings.Split(string(k), "/")
				if len(sp) < 2 {
					fmt.Print()
				}
				l = append(l, PackageListItem{
					Name:    sp[0],
					Version: sp[1],
				})
				return nil
			})
			return nil
		})
		return nil
	})
	return l
}

func (s *BoltStorage) GetParsedPackages() []PackageListItem {
	return s.getPackages("_parsed")
}

func (s *BoltStorage) GetNotParsedPackages() []PackageListItem {
	return s.getPackages("_notparsed")
}

func (s *BoltStorage) GetManualPackages() []PackageListItem {
	// it's without hyphen because manual bucket is just named "manual"
	return s.getPackages("manual")
}
