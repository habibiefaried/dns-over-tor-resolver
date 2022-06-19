package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/habibiefaried/dns-over-tor-resolver/resolvehandler"
	"github.com/miekg/dns"
)

func main() {
	var torResolve *resolvehandler.TorResolve
	resolverbesidetor := getAllBesideTORResolver()

	go func() {
		torResolve = getTORResolver()
	}()

	// attach request handler func
	dns.HandleFunc(".", func(w dns.ResponseWriter, r *dns.Msg) {
		m := new(dns.Msg)
		m.SetReply(r)
		m.Compress = false

		switch r.Opcode {
		case dns.OpcodeQuery:
			for _, q := range m.Question {
				log.Printf("Query for %s\n", q.Name)
				switch q.Qtype {
				case dns.TypeA:
					for _, u := range resolverbesidetor["local"] {
						rr, err := u.Resolve(q.Name) // TODO: to support multiple answer
						if err != nil {
							fmt.Printf("[ERROR]   no answer from local & cache resolver '%v': %v\n", u.GetName(), err)
						} else {
							fmt.Printf("[SUCCESS] got answer from local & cache resolver '%v'\n", u.GetName())
							m.Answer = append(m.Answer, rr)
							break
						}
					}

					// TOR function later here
					if torResolve != nil {
						rr, err := torResolve.Resolve(q.Name)
						if err != nil {
							fmt.Printf("[ERROR]   no answer from main TOR: %v\n", err)
						} else {
							fmt.Printf("[SUCCESS] got answer from main TOR\n")
							m.Answer = append(m.Answer, rr)
							break
						}
					} else {
						fmt.Println("[WARN] TOR is not initialized yet...")
					}

					for _, u := range resolverbesidetor["fallback"] {
						rr, err := u.Resolve(q.Name)
						if err != nil {
							fmt.Printf("[ERROR]   no answer from fallback resolver '%v': %v\n", u.GetName(), err)
						} else {
							fmt.Printf("[SUCCESS] got answer from fallback resolver '%v'\n", u.GetName())
							m.Answer = append(m.Answer, rr)
							break
						}
					}
				}
			}
		}

		w.WriteMsg(m)
	})

	// start server
	port := 5353
	server := &dns.Server{Addr: ":" + strconv.Itoa(port), Net: "udp"}
	log.Printf("Listening at %d/udp\n", port)

	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}
}
