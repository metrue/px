package agent

import (
	"os"
	"testing"

	"github.com/boltdb/bolt"
)

func TestStore(t *testing.T) {
	t.Run("Set/Get/Delete", func(t *testing.T) {
		dbFile := "px.test.db"
		db, err := bolt.Open(dbFile, 0600, nil)
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()
		defer os.Remove(dbFile)

		store := NewStore(db)

		k := "key-a"
		_, err = store.Get(k)
		if err == nil {
			t.Fatalf("should get bucket not existed error")
		}

		expect := "value-a"
		err = store.Set(k, expect)
		if err != nil {
			t.Fatal(err)
		}

		v, err := store.Get(k)
		if err != nil {
			t.Fatal(err)
		}

		if string(v) != expect {
			t.Fatalf("should get %s but got %v", expect, v)
		}

		err = store.Delete(k)
		if err != nil {
			t.Fatal(err)
		}

		v, err = store.Get(k)
		if err != nil {
			t.Fatal(err)
		}
		if string(v) != "" {
			t.Fatal("should get an empty since deleted")
		}

	})
}
