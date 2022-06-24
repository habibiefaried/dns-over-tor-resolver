package cachehandler

import (
	"math/rand"
	"os"
	"testing"
	"time"

	faker "github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
)

func TestSQLite(t *testing.T) {
	sq := SqliteHandler{
		FileName:          "testingcache.sqlite",
		ExpireAgeinSecond: 3600,
	}

	err := sq.Init()
	if err != nil {
		t.Fatal(err)
	}
	defer sq.Close()

	for i := 1; i <= rand.Intn(100); i++ {
		email := faker.Email()
		password := faker.Password()

		err = sq.Put(email, password, "testing", time.Now().Unix())
		if err != nil {
			t.Fatal(err)
		}

		ret, err := sq.Get(email) // take only record that age less than N sec
		if err != nil {
			t.Log(err)
		}

		assert.NotEqual(t, len(ret), 0)
		assert.Equal(t, ret[0], password, "The two words should be the same.")
	}

	oldRecordNum := rand.Intn(100)
	for i := 1; i <= oldRecordNum; i++ {
		email := faker.Email()

		err = sq.Put(email, "NOt imP0rt4nt!", "testing", 1656087061) // old timestamp
		if err != nil {
			t.Fatal(err)
		}
		ret, err := sq.Get(email)
		if err != nil {
			t.Log(err)
		}

		assert.Equal(t, len(ret), 0) // should zero because the timestamp is so old
	}

	cleanupCount, err := sq.CleanUp()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, int64(oldRecordNum), cleanupCount)
	t.Logf("Cleanup %v record test success\n", oldRecordNum)
	// handle multiple values over same key

	for i := 1; i <= rand.Intn(100); i++ {
		email := faker.Email()
		answerQuery := []string{}

		for j := 1; j <= rand.Intn(5); j++ {
			password := faker.Password()
			answerQuery = append(answerQuery, password)
			err = sq.Put(email, password, "testing", time.Now().Unix())
			if err != nil {
				t.Fatal(err)
			}
		}

		ret, err := sq.Get(email)
		if err != nil {
			t.Log(err)
		}

		assert.Equal(t, len(answerQuery), len(ret))

		for _, ans := range answerQuery {
			isStringExist := false
			for _, v := range ret {
				if v == ans {
					isStringExist = true
					break
				}
			}

			if !isStringExist {
				t.Fatalf("the string %v doesn't exist? ", ans)
			}
		}

		t.Logf("key %v has detected %v values\n", email, len(answerQuery))
	}
	t.Log("Multiple record testing done...")

	e := os.Remove("testingcache.sqlite")
	if e != nil {
		t.Fatal(e)
	}

}
