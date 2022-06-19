package resolvehandler

import (
	"testing"

	dotdns "github.com/ncruces/go-dns"
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
