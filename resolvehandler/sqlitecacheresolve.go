package resolvehandler

import (
	"fmt"

	"github.com/habibiefaried/dns-over-tor-resolver/cachehandler"
	"github.com/miekg/dns"
)

type SqliteCacheResolve struct {
	SQLiteHandler cachehandler.SqliteHandler
}

func (s *SqliteCacheResolve) Init() error {
	return nil
}

func (s *SqliteCacheResolve) Resolve(q string) (dns.RR, error) {
	ret, err := s.SQLiteHandler.Get(q)
	if err != nil {
		return nil, err
	}
	return dns.NewRR(fmt.Sprintf("%s A %s", q, *ret))
}

func (s *SqliteCacheResolve) GetName() string {
	return "SQLite-Cache"
}

func (s *SqliteCacheResolve) Close() error {
	fmt.Println("SQLite-cache is closing...")
	return nil
}
