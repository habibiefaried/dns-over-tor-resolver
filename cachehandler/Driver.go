package cachehandler

func InitCachingSystem() (*[]CacheHandler, error) {
	// cache loading
	ret := []CacheHandler{}

	sq := &SqliteHandler{
		FileName:          "dnscache.sqlite",
		ExpireAgeinSecond: 3600,
	}
	err := sq.Init()
	if err != nil {
		return nil, err
	}

	ret = append(ret, sq)
	return &ret, nil
}
