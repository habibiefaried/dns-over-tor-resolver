package resolvehandler

import (
	"fmt"

	"github.com/habibiefaried/dns-over-tor-resolver/cachehandler"
	"github.com/miekg/dns"
)

type CacheResolve struct {
	CacheHandler cachehandler.CacheHandler
}

func (s *CacheResolve) Init() error {
	return nil
}

func (s *CacheResolve) Resolve(q string) (dns.RR, error) {
	ret, err := s.CacheHandler.Get(q)
	if err != nil {
		return nil, err
	}
	return dns.NewRR(fmt.Sprintf("%s 60 IN A %s", q, *ret))
}

func (s *CacheResolve) GetName() string {
	return "Cache"
}

func (s *CacheResolve) Close() error {
	fmt.Println("cache is closing...")
	return nil
}
