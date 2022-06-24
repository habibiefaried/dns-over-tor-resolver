package resolvehandler

import "github.com/miekg/dns"

// ResolveHandler is interface to handle all resolving options, currently we support local memory and tor network
// we can add another options in the future
type ResolveHandler interface {
	Init() error                      // To initialize upstream
	Resolve(string) ([]dns.RR, error) // To resolve the query from DNS to A
	GetName() string
	Close() error
}
