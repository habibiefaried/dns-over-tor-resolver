package resolvehandler

import (
	"fmt"

	"github.com/habibiefaried/dns-over-tor-resolver/cachehandler"
	"github.com/miekg/dns"
)

type CacheResolve struct {
	CacheHandler cachehandler.CacheHandler
	DnsTTL       int
}

func (s *CacheResolve) Init() error {
	if s.DnsTTL == 0 {
		s.DnsTTL = 10 // default value
	}
	return nil
}

func (s *CacheResolve) Resolve(q string) ([]dns.RR, error) {
	retRR := []dns.RR{}

	ret, err := s.CacheHandler.Get(q)
	if err != nil {
		return retRR, err
	}

	for _, relm := range ret {
		d, err := dns.NewRR(fmt.Sprintf("%s %v IN A %s", q, s.DnsTTL, relm))
		if err != nil {
			return nil, err
		}

		retRR = append(retRR, d)
	}

	if len(retRR) == 0 {
		return nil, fmt.Errorf("no record found here")
	} else {
		return retRR, nil
	}
}

func (s *CacheResolve) GetName() string {
	return "Cache"
}

func (s *CacheResolve) Close() error {
	fmt.Println("cache is closing...")
	return nil
}
