package resolvehandler

import (
	"fmt"

	"github.com/miekg/dns"
)

// MemoryResolve implement ResolveHandler so we can
type MemoryResolve struct {
	Records map[string]string
	Name    string
}

func (m *MemoryResolve) GetName() string {
	return m.Name
}

func (m *MemoryResolve) Init() error {
	return nil
}

func (m *MemoryResolve) Resolve(q string) (dns.RR, error) {
	if val, ok := m.Records[q]; ok {
		return dns.NewRR(fmt.Sprintf("%s 60 IN A %s", q, val))
	} else {
		return nil, fmt.Errorf("query not found")
	}
}

func (m *MemoryResolve) Close() error {
	fmt.Println("Memory resolve is closing...")
	return nil
}
