package resolvehandler

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/cretz/bine/tor"
	"github.com/habibiefaried/dns-over-tor-resolver/cachehandler"
	"github.com/miekg/dns"
)

type TorResolve struct {
	OnionDNSServer string
	Name           string
	DNSCache       []cachehandler.CacheHandler // to support cache
	intlresolve    *net.Resolver
	dialCancel     context.CancelFunc
	conn           net.Conn
}

func (tr *TorResolve) Init() error {
	fmt.Println("Starting and registering onion service, please wait a couple of minutes...")
	t, err := tor.Start(context.TODO(), &tor.StartConf{EnableNetwork: true, DebugWriter: nil})
	if err != nil {
		return err
	}
	var dialCtx context.Context
	dialCtx, tr.dialCancel = context.WithTimeout(context.Background(), time.Minute)
	dialer, err := t.Dialer(dialCtx, nil)
	if err != nil {
		return err
	}

	// For testing the DNS resolve via TOR network
	fmt.Println("Testing and building network...")
	tr.intlresolve = &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			tr.conn, err = dialer.Dial("tcp", tr.OnionDNSServer)
			if err != nil {
				log.Println("error dialing remote addr", err)
				return nil, err
			}
			return tr.conn, err
		},
	}
	ip, err := tr.intlresolve.LookupHost(context.Background(), "puredns.org")
	if err != nil {
		return err
	}
	fmt.Println("TOR net resolve success: " + ip[0])

	return nil
}

func (tr *TorResolve) Resolve(q string) ([]dns.RR, error) {
	ret := []dns.RR{}

	ips, err := tr.intlresolve.LookupHost(context.Background(), q)
	if err != nil {
		return nil, err
	} else {
		for _, ip := range ips {
			if net.ParseIP(ip).To4() != nil {

				if tr.DNSCache != nil {
					for _, v := range tr.DNSCache {
						err := v.Put(q, ip, "TOR")
						if err != nil {
							fmt.Printf("Error while putting on cache %v\n", err)
						}
					}
				}

				c, err := dns.NewRR(fmt.Sprintf("%s 60 IN A %s", q, ip))
				if err != nil {
					return nil, fmt.Errorf("got error for creating new record: %v", err)
				}

				ret = append(ret, c)
			}
		}

		return ret, nil
	}
}

func (tr *TorResolve) GetName() string {
	return tr.Name
}

// Close function is a function to close any open connections/processes to upstream
func (tr *TorResolve) Close() error {
	fmt.Println("TOR resolver is closing...")
	err := tr.conn.Close()
	if err != nil {
		return err
	}
	tr.dialCancel()
	return nil
}
