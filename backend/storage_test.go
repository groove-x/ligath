package main

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"go.etcd.io/bbolt"
)

const input1 = `
{
 "host": "yoyoyo",
 "parsed": [
  {
   "name": "haveyoumooedtoday",
   "version": "1.1.1-1",
   "copyrights": [
    {
     "copyright": "2007-2019 Hatsune Miku <miku@example.com>",
     "file_range": [
      "*"
     ],
     "license": {
      "name": "GPL-3+",
      "machine_readable_name": "GPL-3+",
      "body": "This is a test"
     }
    }
   ],
   "raw_copyright": "Original copyright should be here"
  }
 ],
 "not_parsed": [
  {
   "name": "brasero-cdrkit",
   "version": "1.2.3-4",
   "copyrights": [],
   "raw_copyright": ""
  }
 ]
}
`

const input2 = `
{
 "host": "yoyoyo",
 "parsed": [
  {
   "name": "haveyoumooedtoday",
   "version": "1.1.1-2",
   "copyrights": [
    {
     "copyright": "2007-2019 Hatsune Miku <miku@example.com>",
     "file_range": [
      "*"
     ],
     "license": {
      "name": "GPL-3+",
      "machine_readable_name": "GPL-3+",
      "body": "This is a test"
     }
    }
   ],
   "raw_copyright": "Newer Version!!"
  }
 ],
 "not_parsed": [
 ]
}
`

func deleteTestDB(t *testing.T) {
	err := os.Remove("test.db")
	if err != nil {
		t.Fatal("failed to remove test.db:", err)
	}
}

func TestSimpleMigration(t *testing.T) {
	b, err := NewBoltStorage("test.db")
	if err != nil {
		t.Fatal(err)
	}

	defer deleteTestDB(t)
	defer b.Close()

	migrated, err := b.migrate("/hoge/foobar", []byte(input1))
	if err != nil {
		t.Fatal(err)
	} else if !migrated {
		t.Fatal("not migrated")
	}

	err = b.db.View(func(tx *bbolt.Tx) error {
		parsed := tx.Bucket([]byte("foobar_parsed"))
		if parsed == nil {
			t.Fatal("bucket foobar_parsed should exist")
		}

		raw := parsed.Get([]byte("haveyoumooedtoday/1.1.1-1"))
		var res Package
		err = json.Unmarshal(raw, &res)
		if err != nil {
			t.Fatal("failed to unmarshal:", err)
		}

		if res.Name != "haveyoumooedtoday" {
			t.Fatalf("different name: Expected=haveyoumooedtoday Actual=%s", res.Name)
		}

		if res.Version != "1.1.1-1" {
			t.Fatalf("different version: Expected=haveyoumooedtoday Actual=%s", res.Name)
		}

		if res.RawCopyright != "Original copyright should be here" {
			t.Fatalf("different raw copyright: Expected=This is a test Actual=%s", res.RawCopyright)
		}

		cr := res.Copyrights[0]

		if cr.FileRange[0] != "*" {
			t.Fatalf("different file range: Expected=* Actual=%s", cr.FileRange)
		}

		if cr.Copyright != "2007-2019 Hatsune Miku <miku@example.com>" {
			t.Fatalf(
				"different copyright: Expected=2007-2019 Hatsune Miku <miku@example.com> Actual=%s",
				cr.Copyright,
			)
		}

		if cr.License.Name != "GPL-3+" {
			t.Fatalf(
				"different license name: Expected=GPL-3+ Actual=%s",
				cr.License.Name,
			)
		}

		if cr.License.MachineReadableName != "GPL-3+" {
			t.Fatalf(
				"different machine-readable license name: Expected=GPL-3+ Actual=%s",
				cr.License.MachineReadableName,
			)
		}

		if cr.License.Body != "This is a test" {
			t.Fatalf("different license name: Expected=GPL-3+ Actual=%s", cr.License.Name)
		}

		return nil
	})
}

func TestGetPackage(t *testing.T) {
	b, err := NewBoltStorage("test.db")
	if err != nil {
		t.Fatal(err)
	}

	defer deleteTestDB(t)
	defer b.Close()

	migrated, err := b.migrate("/hoge/foo", []byte(input1))
	if err != nil {
		t.Fatal(err)
	} else if !migrated {
		t.Fatal("not migrated")
	}

	pkg, err := b.GetPackage("haveyoumooedtoday", "1.1.1-1")
	if err != nil {
		t.Fatal(err)
	}

	if pkg.Name != "haveyoumooedtoday" {
		t.Fatalf("different name: Expected=haveyoumooedtoday Actual=%s", pkg.Name)
	}

	if pkg.Version != "1.1.1-1" {
		t.Fatalf("different version: Expected=1.1.1-1 Actual=%s", pkg.Version)
	}
}

func TestGetMultipleVersion(t *testing.T) {
	b, err := NewBoltStorage("test.db")
	if err != nil {
		t.Fatal(err)
	}

	defer deleteTestDB(t)
	defer b.Close()

	migrated, err := b.migrate("/hoge/foo", []byte(input1))
	if err != nil {
		t.Fatal(err)
	} else if !migrated {
		t.Fatal("not migrated")
	}

	migrated, err = b.migrate("/hoge/bar", []byte(input2))
	if err != nil {
		t.Fatal(err)
	} else if !migrated {
		t.Fatal("not migrated")
	}

	pkg, err := b.GetPackage("haveyoumooedtoday", "1.1.1-1")
	if err != nil {
		t.Fatal(err)
	}

	if pkg.Version != "1.1.1-1" {
		t.Fatalf("different version: Expected=1.1.1-1 Actual=%s", pkg.Version)
	}

	if pkg.RawCopyright != "Original copyright should be here" {
		t.Fatalf("different copyright: Expected=Original copyright should be here Actual=%s", pkg.RawCopyright)
	}

	pkg, err = b.GetPackage("haveyoumooedtoday", "1.1.1-2")
	if err != nil {
		t.Fatal(err)
	}

	if pkg.Version != "1.1.1-2" {
		t.Fatalf("different version: Expected=1.1.1-2 Actual=%s", pkg.Version)
	}

	if pkg.RawCopyright != "Newer Version!!" {
		t.Fatalf("different copyright: Expected=Newer Version!! Actual=%s", pkg.RawCopyright)
	}
}

func TestIterateUtil(t *testing.T) {
	b, err := NewBoltStorage("test.db")
	if err != nil {
		t.Fatal(err)
	}

	defer deleteTestDB(t)
	defer b.Close()

	migrated, err := b.migrate("/hoge/foo", []byte(input1))
	if err != nil {
		t.Fatal(err)
	} else if !migrated {
		t.Fatal("not migrated")
	}

	migrated, err = b.migrate("/hoge/bar", []byte(input2))
	if err != nil {
		t.Fatal(err)
	} else if !migrated {
		t.Fatal("not migrated")
	}

	// 1. Prefix
	var found bool
	err = IterateBucketsItems(
		b.db,
		&IterateOption{
			BucketPrefix: "fo",
		},
		func(k, v []byte) error {
			if bytes.Equal(k, []byte("haveyoumooedtoday/1.1.1-1")) {
				found = true
			}
			return nil
		},
	)

	if err != nil {
		t.Fatal(err)
	} else if !found {
		t.Fatal("expected key was not found: haveyoumooedtoday/1.1.1-1")
	}

	// 2. Suffix
	found = false
	err = IterateBucketsItems(
		b.db,
		&IterateOption{
			BucketSuffix: "r_parsed",
		},
		func(k, v []byte) error {
			if bytes.Equal(k, []byte("haveyoumooedtoday/1.1.1-2")) {
				found = true
			}
			return nil
		},
	)

	if err != nil {
		t.Fatal(err)
	} else if !found {
		t.Fatal("expected key was not found: haveyoumooedtoday/1.1.1-2")
	}

	// 3. Exact
	found = false
	err = IterateBucketsItems(
		b.db,
		&IterateOption{
			BucketExact: "bar_parsed",
		},
		func(k, v []byte) error {
			if bytes.Equal(k, []byte("haveyoumooedtoday/1.1.1-2")) {
				found = true
			}
			return nil
		},
	)

	if err != nil {
		t.Fatal(err)
	} else if !found {
		t.Fatal("expected key was not found: haveyoumooedtoday/1.1.1-2")
	}

	// 4. Not found
	err = IterateBucketsItems(
		b.db,
		&IterateOption{
			BucketExact: "bucketwhichdoesnotexist",
		},
		func(k, v []byte) error {
			return nil
		},
	)

	// error should happen
	if err == nil {
		t.Fatal("no expected error happend")
	}
}
