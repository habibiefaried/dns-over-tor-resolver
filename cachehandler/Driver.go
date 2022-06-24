package cachehandler

import "github.com/habibiefaried/dns-over-tor-resolver/config"

func InitCachingSystem(c *config.Config) (*[]CacheHandler, error) {
	// cache loading
	ret := []CacheHandler{}

	sq := &SqliteHandler{
		FileName:          "dnscache.sqlite",
		ExpireAgeinSecond: c.CacheTTL,
	}
	err := sq.Init()
	if err != nil {
		return nil, err
	}

	ret = append(ret, sq)
	return &ret, nil
}
