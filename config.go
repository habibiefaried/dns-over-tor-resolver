package main

import (
	"log"

	"github.com/habibiefaried/dns-over-tor-resolver/resolvehandler"
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
				upstreamTOR.Close()
				log.Fatalf("error initializing tor network after %v tries: %v", maxTries, err)
			} else if err != nil {
				trial++
				log.Printf("trial num %v, got error: %v\n", trial, err)
				upstreamTOR.Close()
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

	// ** INIT ALL UPSTREAMS HERE **
	// 1. Local
	upstreamLocal := resolvehandler.MemoryResolve{
		Name:    "Manual",
		Records: c.Manual,
	}
	upstreamLocal.Init()
	defer upstreamLocal.Close()
	upstreams["local"] = append(upstreams["local"], &upstreamLocal)

	return upstreams
}
