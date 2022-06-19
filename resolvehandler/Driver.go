package resolvehandler

import (
	"fmt"
	"log"

	"github.com/habibiefaried/dns-over-tor-resolver/cachehandler"
	"github.com/habibiefaried/dns-over-tor-resolver/config"
	dotdns "github.com/ncruces/go-dns"
)

func GetTORResolver(c *config.Config, ch []cachehandler.CacheHandler) *TorResolve {
	var err error
	maxTries := 10

	if (c.Tor.Address != "") && (c.Tor.Port != "") {
		trial := 1
		for {
			upstreamTOR := TorResolve{
				OnionDNSServer: c.Tor.Address + ":" + c.Tor.Port,
				Name:           "TOR Network",
				DNSCache:       ch,
			}
			err = upstreamTOR.Init()
			if err != nil && trial >= maxTries {
				log.Fatalf("error initializing tor network after %v tries: %v", maxTries, err)
			} else if err != nil {
				trial++
				log.Printf("trial num %v, got error: %v\n", trial, err)
			} else {
				return &upstreamTOR
			}
		}
	}
	return nil
}

func GetAllBesideTORResolver(c *config.Config, ch []cachehandler.CacheHandler) map[string][]ResolveHandler {
	upstreams := make(map[string][]ResolveHandler)
	// ** INIT ALL UPSTREAMS (BESIDE TOR) HERE **
	// 1. Cache
	for _, c := range ch {
		upstreams["local"] = append(upstreams["local"], &CacheResolve{
			CacheHandler: c,
		})
	}

	// 2. Local
	upstreamLocal := MemoryResolve{
		Name:    "Manual",
		Records: c.Manual,
	}
	upstreamLocal.Init()
	upstreams["local"] = append(upstreams["local"], &upstreamLocal)

	// 3. DoT, hardcoded address for now
	dts := []DoTResolve{
		{
			ServerHosts: "cloudflare-dns.com",
			ServerOpts: []dotdns.DoTOption{
				dotdns.DoTAddresses("1.1.1.1", "1.0.0.1", "2606:4700:4700::1111", "2606:4700:4700::1001"),
			},
			DNSCache: ch,
		},
		{
			ServerHosts: "dns.google",
			ServerOpts: []dotdns.DoTOption{
				dotdns.DoTAddresses("8.8.8.8", "8.8.4.4", "2001:4860:4860::8888", "2001:4860:4860::8844"),
			},
			DNSCache: ch,
		},
	}

	for _, v := range dts {
		err := v.Init()
		if err != nil {
			fmt.Printf("Error while initializing DOT %v: %v\n", v.ServerHosts, err)
		} else {
			upstreams["fallback"] = append(upstreams["fallback"], &v)
		}
	}
	return upstreams
}
