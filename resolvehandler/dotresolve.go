package resolvehandler

import (
	"context"
	"fmt"
	"net"

	"github.com/habibiefaried/dns-over-tor-resolver/cachehandler"
	"github.com/miekg/dns"
	dotdns "github.com/ncruces/go-dns"
)

type DoTResolve struct {
	ServerHosts string
	ServerOpts  []dotdns.DoTOption
	DNSCache    []cachehandler.CacheHandler // to support cache
	resolver    *net.Resolver
}

func (dt *DoTResolve) Init() error {
	var err error
	dt.resolver, err = dotdns.NewDoTResolver(dt.ServerHosts, dt.ServerOpts...)
	if err != nil {
		return fmt.Errorf("newDoTResolver %v error = %v", dt.ServerHosts, err)
	}
	return nil
}

func (dt *DoTResolve) Resolve(q string) (dns.RR, error) {
	ips, err := dt.resolver.LookupIPAddr(context.TODO(), q)
	if err != nil {
		return nil, err
	}

	for _, ip := range ips {
		if net.ParseIP(ip.String()) != nil {
			for _, v := range dt.DNSCache {
				err := v.Put(q, ip.String(), fmt.Sprintf("DOT-%v", dt.ServerHosts))
				if err != nil {
					fmt.Printf("Error while putting on cache %v\n", err)
				}
			}

			return dns.NewRR(fmt.Sprintf("%s 60 IN A %s", q, ip.String())) // TODO: return multiple value
		}
	}

	return nil, fmt.Errorf("all of these IPs not valid for IPv4 format: %v", ips)
}

func (dt *DoTResolve) GetName() string {
	return fmt.Sprintf("DOT-%v", dt.ServerHosts)
}

func (dt *DoTResolve) Close() error {
	fmt.Println("DOT resolver is closing...")
	return nil
}
