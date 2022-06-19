package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/miekg/dns"
)

func main() {
	resolverbesidetor := getAllBesideTORResolver()

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
							fmt.Printf("[ERROR]   no answer from resolver '%v': %v\n", u.GetName(), err)
						} else {
							fmt.Printf("[SUCCESS] got answer from resolver '%v'\n", u.GetName())
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
