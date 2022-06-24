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

func (s *CacheResolve) Resolve(q string) ([]dns.RR, error) {
	retRR := []dns.RR{}

	ret, err := s.CacheHandler.Get(q)
	if err != nil {
		return retRR, err
	}

	for _, relm := range ret {
		d, err := dns.NewRR(fmt.Sprintf("%s 60 IN A %s", q, relm))
		if err != nil {
			return nil, err
		}

		retRR = append(retRR, d)
	}

	return retRR, nil
}

func (s *CacheResolve) GetName() string {
	return "Cache"
}

func (s *CacheResolve) Close() error {
	fmt.Println("cache is closing...")
	return nil
}
