package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/habibiefaried/dns-over-tor-resolver/resolvehandler"
	"github.com/miekg/dns"
)

var upstreams []resolvehandler.ResolveHandler

func parseAnswerQuery(m *dns.Msg) {
	for _, q := range m.Question {
		log.Printf("Query for %s\n", q.Name)
		switch q.Qtype {
		case dns.TypeA:
			for _, u := range upstreams {
				rr, err := u.Resolve(q.Name)
				if err != nil {
					fmt.Printf("[ERROR from '%v'] %v\n", u.GetName(), err)
				} else {
					fmt.Printf("Got answer from resolver '%v'!\n", u.GetName())
					m.Answer = append(m.Answer, rr)
					break
				}
			}
		}
	}
}

func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		parseAnswerQuery(m)
	}

	w.WriteMsg(m)
}

func main() {
	var err error
	upstreams = []resolvehandler.ResolveHandler{}

	upstreamLocal := resolvehandler.MemoryResolve{
		Name: "Test",
	}
	err = upstreamLocal.Init()
	if err != nil {
		panic(err)
	}
	defer upstreamLocal.Close()
	upstreams = append(upstreams, &upstreamLocal)

	upstreamTOR := resolvehandler.TorResolve{
		OnionDNSServer: "pdnsc3ocqx2qj2irfeinfduqckjyjsjvruyfmbt2i5eszy3k2u4gcoqd.onion:53",
		Name:           "TOR Network",
	}
	err = upstreamTOR.Init()
	if err != nil {
		panic(err)
	}
	defer upstreamTOR.Close()
	upstreams = append(upstreams, &upstreamTOR)

	// attach request handler func
	dns.HandleFunc(".", handleDnsRequest)
	// start server
	port := 5353
	server := &dns.Server{Addr: ":" + strconv.Itoa(port), Net: "udp"}
	log.Printf("Listening at %d/udp\n", port)

	err = server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}
}
