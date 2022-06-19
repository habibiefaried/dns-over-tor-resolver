package resolvehandler

import (
	"context"
	"fmt"
	"net"

	"github.com/miekg/dns"
	dotdns "github.com/ncruces/go-dns"
)

type DoTResolve struct {
	ServerHosts string
	ServerOpts  []dotdns.DoTOption
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

	return dns.NewRR(fmt.Sprintf("%s A %s", q, ips[0].String())) // TODO: return multiple value
}

func (dt *DoTResolve) GetName() string {
	return fmt.Sprintf("DOT-%v", dt.ServerHosts)
}

func (dt *DoTResolve) Close() error {
	return nil
}
