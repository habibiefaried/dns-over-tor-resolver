package cachehandler

type CacheHandler interface {
	Init() error
	Put(string, string, string) error // set Key value
	Get(string) (*string, error)      // get value
	Close() error
}
