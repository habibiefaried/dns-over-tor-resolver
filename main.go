package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/habibiefaried/dns-over-tor-resolver/cachehandler"
	"github.com/habibiefaried/dns-over-tor-resolver/config"
	"github.com/habibiefaried/dns-over-tor-resolver/resolvehandler"
	"github.com/miekg/dns"
)

func main() {
	c, err := config.ReadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	caches, err := cachehandler.InitCachingSystem(c)
	if err != nil {
		log.Fatal(err)
	}

	resolvers := resolvehandler.GetAllBesideTORResolver(c, *caches)
	resolvers["tor"] = nil
	go func() {
		resolvers["tor"] = append(resolvers["tor"], resolvehandler.GetTORResolver(c, *caches))
	}()
	keysInSorted := []string{"local", "tor", "fallback"}

	// attach request handler func
	dns.HandleFunc(".", func(w dns.ResponseWriter, r *dns.Msg) {
		m := new(dns.Msg)
		m.SetReply(r)
		m.Compress = false

	OuterLoop:
		switch r.Opcode {
		case dns.OpcodeQuery:
			for _, q := range m.Question {
				log.Printf("Query for %s\n", q.Name)
				switch q.Qtype {
				case dns.TypeA:
					for _, sKey := range keysInSorted {
						if resolvers[sKey] != nil {
							for _, u := range resolvers[sKey] {
								rr, err := u.Resolve(q.Name)
								if err != nil {
									fmt.Printf("[ERROR]   no answer from resolver '%v': %v\n", u.GetName(), err)
								} else {
									fmt.Printf("[SUCCESS] got answer from resolver '%v'\n", u.GetName())
									m.Answer = append(m.Answer, rr...)
									break OuterLoop
								}
							}
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

	err = server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}
}
