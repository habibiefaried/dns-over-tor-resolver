package resolvehandler

import (
	"testing"

	"github.com/habibiefaried/dns-over-tor-resolver/cachehandler"
	dotdns "github.com/ncruces/go-dns"
	"github.com/stretchr/testify/assert"
)

func TestDOT(t *testing.T) {
	dts := []DoTResolve{
		{
			ServerHosts: "cloudflare-dns.com",
			ServerOpts: []dotdns.DoTOption{
				dotdns.DoTAddresses("1.1.1.1", "1.0.0.1", "2606:4700:4700::1111", "2606:4700:4700::1001"),
			},
		},
		{
			ServerHosts: "dns.google",
			ServerOpts: []dotdns.DoTOption{
				dotdns.DoTAddresses("8.8.8.8", "8.8.4.4", "2001:4860:4860::8888", "2001:4860:4860::8844"),
			},
		},
	}

	for _, v := range dts {
		err := v.Init()
		if err != nil {
			t.Fatal(err)
		}

		res, err := v.Resolve("puredns.org")
		if err != nil {
			t.Fatal(err)
		}

		t.Log(res.String())

		_, err = v.Resolve("nonexistenthosts.pvt")
		t.Log(err)
	}
}

func TestSQLiteResolve(t *testing.T) {
	sq := &cachehandler.SqliteHandler{
		FileName: "testingcache.sqlite",
	}

	err := sq.Init()
	if err != nil {
		t.Fatal(err)
	}
	defer sq.Close()

	err = sq.Put("google.com.", "8.8.8.8", "testing")
	if err != nil {
		t.Fatal(err)
	}

	sqresolv := SqliteCacheResolve{
		SQLiteHandler: sq,
	}

	rr, err := sqresolv.Resolve("google.com.")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "google.com.	3600	IN	A	8.8.8.8", rr.String())
}
