package main

import (
	"fmt"

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

func applyConfig() []resolvehandler.ResolveHandler {
	upstreams := []resolvehandler.ResolveHandler{}
	var err error
	c, err := readConfig()
	if err != nil {
		panic(err)
	}

	// ** INIT ALL UPSTREAMS HERE **
	// 1. Local
	upstreamLocal := resolvehandler.MemoryResolve{
		Name:    "Manual",
		Records: c.Manual,
	}
	upstreamLocal.Init()
	defer upstreamLocal.Close()
	upstreams = append(upstreams, &upstreamLocal)

	// 2. TOR
	if (c.Tor.Address != "") && (c.Tor.Port != "") {
		upstreamTOR := resolvehandler.TorResolve{
			OnionDNSServer: c.Tor.Address + ":" + c.Tor.Port,
			Name:           "TOR Network",
		}
		err = upstreamTOR.Init()
		if err != nil {
			fmt.Println(err)
		} else {
			upstreams = append(upstreams, &upstreamTOR)
		}
		defer upstreamTOR.Close()
	}

	return upstreams
}
