package cachehandler

import (
	"os"
	"testing"

	faker "github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
)

func TestSQLite(t *testing.T) {
	sq := SqliteHandler{
		FileName: "testingcache.sqlite",
	}

	err := sq.Init()
	if err != nil {
		t.Fatal(err)
	}
	defer sq.Close()

	for i := 1; i <= 100; i++ {
		email := faker.Email()
		password := faker.Password()

		err = sq.Put(email, password, "testing")
		if err != nil {
			t.Fatal(err)
		}

		ret, err := sq.Get(email)
		if err != nil {
			t.Log(err)
		}
		assert.Equal(t, *ret, password, "The two words should be the same.")
	}

	for i := 1; i <= 50; i++ {
		// Check double entry
		email := faker.FirstName()
		password := faker.Password()
		password2 := faker.Password()

		err = sq.Put(email, password, "testing")
		if err != nil {
			t.Fatal(err)
		}

		err = sq.Put(email, password2, "testing")
		if err != nil {
			t.Fatal(err)
		}

		ret, err := sq.Get(email)
		if err != nil {
			t.Log(err)
		}
		assert.Equal(t, *ret, password2, "The two words should be the same.")
	}

	e := os.Remove("testingcache.sqlite")
	if e != nil {
		t.Fatal(e)
	}

}
