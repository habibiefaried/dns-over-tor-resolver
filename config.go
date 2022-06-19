package main

import (
	"fmt"
	"log"

	"github.com/habibiefaried/dns-over-tor-resolver/resolvehandler"
	dotdns "github.com/ncruces/go-dns"
	"github.com/spf13/viper"
)

type TorConfig struct {
	Address string `mapstructure:"address"`
	Port    string `mapstructure:"port"`
}

type Config struct {
	Tor    TorConfig         `mapstructure:"tor"`
	Manual map[string]string `mapstructure:"manual"`
}

func readConfig() (*Config, error) {
	var config Config
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func getTORResolver() *resolvehandler.TorResolve {
	var err error
	c, err := readConfig()
	if err != nil {
		log.Fatal(err)
	}

	maxTries := 10

	if (c.Tor.Address != "") && (c.Tor.Port != "") {
		trial := 1
		for {
			upstreamTOR := resolvehandler.TorResolve{
				OnionDNSServer: c.Tor.Address + ":" + c.Tor.Port,
				Name:           "TOR Network",
			}
			err = upstreamTOR.Init()
			if err != nil && trial >= maxTries {
				log.Fatalf("error initializing tor network after %v tries: %v", maxTries, err)
			} else if err != nil {
				trial++
				log.Printf("trial num %v, got error: %v\n", trial, err)
			} else {
				return &upstreamTOR
			}
		}
	}
	return nil
}

func getAllBesideTORResolver() map[string][]resolvehandler.ResolveHandler {
	upstreams := make(map[string][]resolvehandler.ResolveHandler)
	var err error
	c, err := readConfig()
	if err != nil {
		log.Fatal(err)
	}

	// ** INIT ALL UPSTREAMS (BESIDE TOR) HERE **
	// 1. Local
	upstreamLocal := resolvehandler.MemoryResolve{
		Name:    "Manual",
		Records: c.Manual,
	}
	upstreamLocal.Init()
	upstreams["local"] = append(upstreams["local"], &upstreamLocal)

	// 2. DoT, hardcoded address for now
	dts := []resolvehandler.DoTResolve{
		{
			ServerHosts: "cloudflare-dns.com",
			ServerOpts: []dotdns.DoTOption{
				dotdns.DoTAddresses("1.1.1.1", "1.0.0.1", "2606:4700:4700::1111", "2606:4700:4700::1001"),
			},
		},
		{
			ServerHosts: "dns.google",
			ServerOpts: []dotdns.DoTOption{
				dotdns.DoTAddresses("8.8.8.8", "8.8.4.4", "2001:4860:4860::8888", "2001:4860:4860::8844"),
			},
		},
	}

	for _, v := range dts {
		err := v.Init()
		if err != nil {
			fmt.Printf("Error while initializing DOT %v: %v\n", v.ServerHosts, err)
		} else {
			upstreams["fallback"] = append(upstreams["fallback"], &v)
		}
	}
	return upstreams
}
