package main

import (
	"bytes"
	"fmt"

	"go.etcd.io/bbolt"
)

type IterateOption struct {
	BucketPrefix string
	BucketSuffix string
	BucketExact  string
}

func IterateBuckets(b *bbolt.DB, option *IterateOption, fn func(b *bbolt.Bucket) error) error {
	var prefix, suffix, exact []byte
	if option != nil {
		if option.BucketPrefix != "" {
			prefix = []byte(option.BucketPrefix)
		}
		if option.BucketSuffix != "" {
			suffix = []byte(option.BucketSuffix)
		}
		if option.BucketExact != "" {
			exact = []byte(option.BucketExact)
		}
	}

	var found bool

	err := b.View(func(tx *bbolt.Tx) error {
		tx.ForEach(func(name []byte, b *bbolt.Bucket) error {
			if exact != nil && !bytes.Equal(name, exact) {
				return nil
			}
			if prefix != nil && !bytes.HasPrefix(name, prefix) {
				return nil
			}
			if suffix != nil && !bytes.HasSuffix(name, suffix) {
				return nil
			}
			found = true
			return fn(b)
		})
		return nil
	})

	if err != nil {
		if option == nil {
			err = fmt.Errorf("failed to iterate over db: %v", err)
		} else {
			err = fmt.Errorf("failed to iterate over db with option: %+v, err: %v", option, err)
		}
	} else if !found {
		err = fmt.Errorf("failed to find buckets with option: %+v", option)
	}

	return err
}

func IterateBucketsItems(b *bbolt.DB, option *IterateOption, fn func(k, v []byte) error) error {
	var prefix, suffix, exact []byte
	if option != nil {
		if option.BucketPrefix != "" {
			prefix = []byte(option.BucketPrefix)
		}
		if option.BucketSuffix != "" {
			suffix = []byte(option.BucketSuffix)
		}
		if option.BucketExact != "" {
			exact = []byte(option.BucketExact)
		}
	}

	var found bool

	err := b.View(func(tx *bbolt.Tx) error {
		tx.ForEach(func(name []byte, b *bbolt.Bucket) error {
			if exact != nil && !bytes.Equal(name, exact) {
				return nil
			}
			if prefix != nil && !bytes.HasPrefix(name, prefix) {
				return nil
			}
			if suffix != nil && !bytes.HasSuffix(name, suffix) {
				return nil
			}
			found = true
			return b.ForEach(fn)
		})
		return nil
	})

	if err != nil {
		if option == nil {
			err = fmt.Errorf("failed to iterate over db: %v", err)
		} else {
			err = fmt.Errorf("failed to iterate over db with option: %+v, err: %v", option, err)
		}
	} else if !found {
		err = fmt.Errorf("failed to find buckets with option: %+v", option)
	}

	return err
}
