package resolvehandler

import (
	"fmt"

	"github.com/miekg/dns"
)

// MemoryResolve implement ResolveHandler so we can
type MemoryResolve struct {
	Records map[string]string
	Name    string
	DnsTTL  int
}

func (m *MemoryResolve) GetName() string {
	return m.Name
}

func (m *MemoryResolve) Init() error {
	if m.DnsTTL == 0 {
		m.DnsTTL = 10 // default value
	}

	return nil
}

func (m *MemoryResolve) Resolve(q string) ([]dns.RR, error) {
	ret := []dns.RR{}

	if val, ok := m.Records[q]; ok {
		c, err := dns.NewRR(fmt.Sprintf("%s %v IN A %s", q, m.DnsTTL, val))
		if err != nil {
			return nil, fmt.Errorf("not correct IP")
		}

		ret = append(ret, c)
		return ret, nil
	} else {
		return nil, fmt.Errorf("query not found")
	}
}

func (m *MemoryResolve) Close() error {
	fmt.Println("Memory resolve is closing...")
	return nil
}
