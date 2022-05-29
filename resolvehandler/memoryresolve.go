package resolvehandler

import (
	"fmt"

	"github.com/miekg/dns"
)

// MemoryResolve implement ResolveHandler so we can
type MemoryResolve struct {
	records map[string]string
	Name    string
}

func (m *MemoryResolve) GetName() string {
	return m.Name
}

func (m *MemoryResolve) Init() error {
	m.records = map[string]string{
		"dns.testing-only.":  "127.0.53.2",
		"local.testing.pvt.": "127.0.4.2",
	}

	return nil
}

func (m *MemoryResolve) Resolve(q string) (dns.RR, error) {
	if val, ok := m.records[q]; ok {
		return dns.NewRR(fmt.Sprintf("%s A %s", q, val))
	} else {
		return nil, fmt.Errorf("query not found")
	}
}

func (m *MemoryResolve) Close() error {
	return nil
}
